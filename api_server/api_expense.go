package main

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	types "github.com/snakehunterr/hacs_app/db_types"
	api_errors "github.com/snakehunterr/hacs_app/db_types/errors"
	validators "github.com/snakehunterr/hacs_app/db_types/validators"
)

func expenseScanRow(e *types.Expense, row *sql.Row) error {
	return row.Scan(&e.ID, &e.Date, &e.Amount, &e.LastEdited)
}

func expenseScanRows(es *[]types.Expense, rows *sql.Rows) error {
	if es == nil {
		return errors.New("*[]types.Expense is nil")
	}

	var _es []types.Expense
	for rows.Next() {
		var e types.Expense

		if err := rows.Scan(&e.ID, &e.Date, &e.Amount, &e.LastEdited); err != nil {
			return err
		}

		_es = append(_es, e)
	}

	*es = _es
	return nil
}

//go:embed sql/expense/expense_get_all.sql
var SQLExpenseGetAllQuery string

// ExpenseAll godoc
// @Summary Get all expenses
// @Schemes http
// @Description Get all expenses
// @Tags expense
// @Produce json
// @Success 200 {array} types.Expense "ok"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /expense/all [get]
func RouteExpenseGetAll(g *gin.Context) {
	var es []types.Expense

	code, err := queryRows(&es, expenseScanRows, SQLExpenseGetAllQuery)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(es) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, es)
}

//go:embed sql/expense/expense_get_by_expense_id.sql
var SQLExpenseGetByIDQuery string

// ExpenseByID godoc
// @Summary Get expense by expense_id
// @Schemes http
// @Description Get expense by expense_id
// @Param id path int true "Expense ID"
// @Tags expense
// @Produce json
// @Success 200 {object} types.Expense "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /expense/id/{id} [get]
func RouteExpenseGetByExpenseID(g *gin.Context) {
	id := g.Param("id")
	if _, err := validators.Int64("id", id, false); err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	var e types.Expense
	code, err := queryRow(&e, expenseScanRow, SQLExpenseGetByIDQuery, id)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	g.JSON(http.StatusOK, e)
}

//go:embed sql/expense/expense_get_by_date.sql
var SQLExpenseGetByDateQuery string

// ExpenseByDate godoc
// @Summary Get expenses by expense_date
// @Schemes http
// @Description Get expenses by expense_date
// @Param date path string true "Expense date"
// @Tags expense
// @Produce json
// @Success 200 {array} types.Expense "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /expense/date/{date} [get]
func RouteExpenseGetByDate(g *gin.Context) {
	date, err := validators.Date("date", g.Param("date"), false)
	if err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	var es []types.Expense
	code, err := queryRows(
		&es, expenseScanRows,
		SQLExpenseGetByDateQuery,
		date.Year(), date.Month(), date.Day(),
	)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(es) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, es)
}

// ExpenseByDateRange godoc
// @Summary Get expenses by date range
// @Schemes http
// @Description Get expenses by date range
// @Param date_start formData string true "Expense date start"
// @Param date_end formData string true "Expense date end"
// @Tags expense
// @Produce json
// @Success 200 {array} types.Expense "ok"
// @Failure 400 {objcet} types.APIResponse "Incorrect parameter"
// @Failure 404 {objcet} types.APIResponse "No rows"
// @Failure 500 {objcet} types.APIResponse "Internal server error"
// @Router /expense/date/range [post]
func RouteExpenseGetByDateRange(g *gin.Context) {
	var (
		apierr     *api_errors.APIError
		date_start time.Time
		date_end   time.Time
	)

	date_start, apierr = validators.Date("date_start", g.PostForm("date_start"), true)
	if apierr != nil {
		goto skip
	}

	date_end, apierr = validators.Date("date_end", g.PostForm("date_end"), true)
	if apierr != nil {
		goto skip
	}

skip:
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: apierr})
		return
	}

	var es []types.Expense
	code, err := queryRows(&es, expenseScanRows, SQLExpenseGetAllQuery)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	var _es []types.Expense
	for _, e := range es {
		if !e.Date.After(date_start) && !e.Date.Equal(date_start) {
			continue
		}
		if !date_end.After(e.Date) && !e.Date.Equal(date_end) {
			continue
		}
		_es = append(_es, e)
	}

	if len(_es) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, _es)
}

//go:embed sql/expense/expense_insert.sql
var SQLExpensePostCreateQuery string

// ExpenseCreate godoc
// @Summary Create new expense
// @Schemes http
// @Description Create new expense
// @Tags expense
// @Param expense_date formData string false "Expense date"
// @Param expense_amount formData number true "Expense amount"
// @Produce json
// @Success 201 {object} types.Expense "New expense"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /expense/new [post]
func RouteExpensePostCreate(g *gin.Context) {
	var (
		apierr         *api_errors.APIError
		expense_date   time.Time
		expense_amount float64
	)

	expense_date, apierr = validators.Date("expense_date", g.PostForm("expense_date"), true)
	if apierr != nil {
		goto skip
	}

	expense_amount, apierr = validators.Float64("expense_amount", g.PostForm("expense_amount"), true)

skip:
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: apierr})
		return
	}

	res, err := db.Exec(SQLExpensePostCreateQuery, expense_date, expense_amount)
	if err != nil {
		logError("db.Exec():", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		logError("res.LastInsertId():", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	e := types.Expense{
		ID:     id,
		Date:   expense_date,
		Amount: expense_amount,
	}

	logInfo(fmt.Sprintf("Created new expense: %#v", e))
	g.JSON(http.StatusCreated, e)
}

//go:embed sql/expense/expense_delete.sql
var SQLExpenseDeleteQuery string

// ExpenseDelete godoc
// @Summary Delete expense by expense_id
// @Schemes http
// @Description Delete expense by expense_id
// @Tags expense
// @Param id path int true "Expense ID"
// @Produce json
// @Success 200 {object} types.APIResponse "Deleted"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /expense/id/{id} [delete]
func RouteExpenseDelete(g *gin.Context) {
	id := g.Param("id")
	if _, err := validators.Int64("id", id, false); err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	if _, err := db.Exec(SQLExpenseDeleteQuery, id); err != nil {
		logError("db.Exec():", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	logInfo("Deleted record expense with expense_id: ", id)
	g.JSON(http.StatusOK, types.APIResponse{Message: "ok"})
}

//go:embed sql/expense/expense_patch.sql
var SQLExpensePatchQuery string

// ExpensePatch godoc
// @Summary Patch expense
// @Schemes http
// @Description Patch expense by expense_id
// @Tags expense
// @Param id path int true "Expense ID"
// @Param expense_date formData string false "Date 'yyyy-mm-dd hh:mm:ss'"
// @Param expense_amount formData number false "Amount"
// @Produce json
// @Success 200 {object} types.APIResponse "Updated"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /expense/id/{id} [patch]
func RouteExpensePatch(g *gin.Context) {
	var (
		apierr     *api_errors.APIError
		expense_id int64
		temp       string
		cache      = map[string]any{}
	)

	expense_id, apierr = validators.Int64("id", g.Param("id"), false)
	if apierr != nil {
		goto skip
	}

	temp = g.PostForm("expense_date")
	if temp != "" {
		var date time.Time

		date, apierr = validators.Date("expense_date", temp, false)
		if apierr != nil {
			goto skip
		}
		cache["expense_date"] = date.Format(validators.DATE_FORMAT)
	}

	temp = g.PostForm("expense_amount")
	if temp != "" {
		var amount float64

		amount, apierr = validators.Float64("expense_amount", temp, false)
		if apierr != nil {
			goto skip
		}
		cache["expense_amount"] = amount
	}

skip:
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: apierr})
		return
	}

	var e types.Expense
	code, apierr := queryRow(&e, expenseScanRow, SQLExpenseGetByIDQuery, expense_id)
	if apierr != nil {
		g.JSON(code, types.APIResponse{Error: apierr})
		return
	}

	if _, ok := cache["expense_date"]; !ok {
		cache["expense_date"] = e.Date.Format(validators.DATE_FORMAT)
	}
	if _, ok := cache["expense_amount"]; !ok {
		cache["expense_amount"] = e.Amount
	}

	_, err := db.Exec(
		SQLExpensePatchQuery,
		cache["expense_date"],
		cache["expense_amount"],
		expense_id,
	)
	if err != nil {
		logError("db.Exec():", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	g.JSON(http.StatusOK, types.APIResponse{Message: "ok"})
}

func init() {
	r := api.Group("/expense")

	r.GET("/all", RouteExpenseGetAll)
	r.GET("/date/:date", RouteExpenseGetByDate)
	r.POST("/date/range", RouteExpenseGetByDateRange)
	r.POST("/new", RouteExpensePostCreate)
	r.GET("/id/:id", RouteExpenseGetByExpenseID)
	r.DELETE("/id/:id", RouteExpenseDelete)
	r.PATCH("/id/:id", RouteExpensePatch)
}
