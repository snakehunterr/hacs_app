package main

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	types "github.com/snakehunterr/hacs_db_types"
	api_errors "github.com/snakehunterr/hacs_db_types/errors"
	validators "github.com/snakehunterr/hacs_db_types/validators"
)

func paymentScanRow(p *types.Payment, row *sql.Row) error {
	return row.Scan(&p.ID, &p.ClientID, &p.RoomID, &p.Date, &p.Amount, &p.LastEdited)
}

func paymentScanRows(ps *[]types.Payment, rows *sql.Rows) error {
	if ps == nil {
		return errors.New("*[]types.Payment is nil")
	}

	_ps := *ps
	for rows.Next() {
		var p types.Payment

		if err := rows.Scan(&p.ID, &p.ClientID, &p.RoomID, &p.Date, &p.Amount, &p.LastEdited); err != nil {
			return err
		}

		_ps = append(_ps, p)
	}

	*ps = _ps
	return nil
}

//go:embed sql/payment/payment_get_all.sql
var SQLPaymentGetAllQuery string

// PaymentAll godoc
// @Summary Get all payments
// @Schemes http
// @Description Get all payments
// @Tags payment
// @Produce json
// @Success 200 {array} types.Payment "ok"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /payment/all [get]
func RoutePaymentGetAll(g *gin.Context) {
	var ps []types.Payment

	code, err := queryRows(&ps, paymentScanRows, SQLPaymentGetAllQuery)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(ps) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, ps)
}

//go:embed sql/payment/payment_get_by_client_id.sql
var SQLPaymentGetByClientIDQuery string

// PaymentAllByClientID godoc
// @Summary Get all payments by client_id
// @Schemes http
// @Description Get all payments by client_id
// @Param id path int true "Client ID"
// @Tags payment
// @Produce json
// @Success 200 {array} types.Payment "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /payment/client/id/{id} [get]
func RoutePaymentGetAllByClientID(g *gin.Context) {
	id := g.Param("id")
	if _, err := validators.Int64("id", id, false); err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	var ps []types.Payment
	code, err := queryRows(&ps, paymentScanRows, SQLPaymentGetByClientIDQuery, id)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(ps) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, ps)
}

//go:embed sql/payment/payment_get_by_room_id.sql
var SQLPaymentGetByRoomIDQuery string

// PaymentGetAllByRoomID godoc
// @Summary Get all payments by room_id
// @Schemes http
// @Description Get all payments by room_id
// @Param id path int true "Room ID"
// @Tags payment
// @Produce json
// @Success 200 {array} types.Payment "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /payment/room/id/{id} [get]
func RoutePaymentGetAllByRoomID(g *gin.Context) {
	id := g.Param("id")
	if _, err := validators.Int64("id", id, false); err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	var ps []types.Payment
	code, err := queryRows(&ps, paymentScanRows, SQLPaymentGetByRoomIDQuery, id)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(ps) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, ps)
}

//go:embed sql/payment/payment_get_by_id.sql
var SQLPaymentGetByPaymentIDQuery string

// PaymentGetByID godoc
// @Summary Get payment by payment_id
// @Schemes http
// @Description Get payment by payment_id
// @Param id path int true "Payment ID"
// @Tags payment
// @Produce json
// @Success 200 {object} types.Payment "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /payment/id/{id} [get]
func RoutePaymentGetByID(g *gin.Context) {
	id := g.Param("id")
	if _, err := validators.Int64("id", id, false); err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	var p types.Payment
	code, err := queryRow(&p, paymentScanRow, SQLPaymentGetByPaymentIDQuery, id)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	g.JSON(http.StatusOK, p)
}

//go:embed sql/payment/payment_get_by_date.sql
var SQLPaymentGetByDateQuery string

// PaymentByDate godoc
// @Summary Get payment by date
// @Schemes http
// @Description Get payment by date
// @Param date path string true "Date 'yyyy-mm-dd hh:mm:ss'"
// @Tags payment
// @Produce json
// @Success 200 {array} types.Payment "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /payment/date/{date} [get]
func RoutePaymentGetByDate(g *gin.Context) {
	date, err := validators.Date("date", g.Param("date"), false)
	if err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	var ps []types.Payment
	code, err := queryRows(&ps, paymentScanRows, SQLPaymentGetByDateQuery, date.Year(), date.Month(), date.Day())
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(ps) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, ps)
}

// PaymentByDateRange godoc
// @Summary Get payment by date range
// @Schemes http
// @Description Get payment by date range
// @Param date_start formData string true "Date 'yyyy-mm-dd hh:mm:ss'"
// @Param date_end formData string true "Date 'yyyy-mm-dd hh:mm:ss'"
// @Tags payment
// @Produce json
// @Success 200 {array} types.Payment "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /payment/date/range [post]
func RoutePaymentGetByDateRange(g *gin.Context) {
	var (
		date_start time.Time
		date_end   time.Time
		err        *api_errors.APIError
	)
	date_start, err = validators.Date("date_start", g.PostForm("date_start"), true)
	if err != nil {
		goto skip
	}

	date_end, err = validators.Date("date_end", g.PostForm("date_end"), true)
	if err != nil {
		goto skip
	}

skip:
	if err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	var ps []types.Payment
	code, err := queryRows(
		&ps, paymentScanRows,
		SQLPaymentGetAllQuery,
	)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	var _ps []types.Payment
	for _, p := range ps {
		if !p.Date.After(date_start) && !p.Date.Equal(date_start) {
			continue
		}
		if !date_end.After(p.Date) && !p.Date.Equal(date_end) {
			continue
		}
		_ps = append(_ps, p)
	}

	if len(_ps) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, _ps)
}

//go:embed sql/payment/payment_insert.sql
var SQLPaymentPostCreateQuery string

// PaymentCreate godoc
// @Summary Create new payment
// @Schemes http
// @Description Create new payment
// @Param client_id formData int true "Client ID"
// @Param room_id formData int true "Room ID"
// @Param payment_date formData string false "Date 'yyyy-mm-dd hh:mm:ss'"
// @Param payment_amount formData number true "Amount"
// @Tags payment
// @Produce json
// @Success 201 {object} types.Payment "New payment"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /payment/new [post]
func RoutePaymentPostCreate(g *gin.Context) {
	var (
		apierr         *api_errors.APIError
		client_id      int64
		room_id        int64
		payment_date   time.Time
		payment_amount float64
	)

	client_id, apierr = validators.Int64("client_id", g.PostForm("client_id"), true)
	if apierr != nil {
		goto skip
	}

	room_id, apierr = validators.Int64("room_id", g.PostForm("room_id"), true)
	if apierr != nil {
		goto skip
	}

	payment_date, apierr = validators.Date("payment_date", g.PostForm("payment_date"), true)
	if apierr != nil {
		goto skip
	}

	payment_amount, apierr = validators.Float64("payment_amount", g.PostForm("payment_amount"), true)

skip:
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: apierr})
		return
	}

	res, err := db.Exec(SQLPaymentPostCreateQuery, client_id, room_id, payment_date, payment_amount)
	if err != nil {
		log.Println("[ERROR] RoutePaymentPostCreate db.Exec():", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	payment_id, err := res.LastInsertId()
	if err != nil {
		logError("sql.Result LastInsertId(): ", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	p := types.Payment{
		ID:       payment_id,
		ClientID: client_id,
		RoomID:   room_id,
		Date:     payment_date,
		Amount:   payment_amount,
	}

	logInfo(fmt.Sprintf("Created new payment: %#v", p))
	g.JSON(http.StatusCreated, p)
}

//go:embed sql/payment/payment_delete.sql
var SQLPaymentDeleteQuery string

// PaymentDelete godoc
// @Summary Delete payment
// @Schemes http
// @Description Delete payment by payment_id
// @Param id path int true "Payment ID"
// @Tags payment
// @Produce json
// @Success 200 {object} types.APIResponse "Deleted"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /payment/id/{id} [delete]
func RoutePaymentDelete(g *gin.Context) {
	id := g.Param("id")
	if _, err := validators.Int64("id", id, false); err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	_, err := db.Exec(SQLPaymentDeleteQuery, id)
	if err != nil {
		logError("db.Exec():", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	logInfo("Deleted record payment with payment_id: ", id)
	g.JSON(http.StatusOK, types.APIResponse{Message: "ok"})
}

//go:embed sql/payment/payment_patch.sql
var SQLPaymentPatchQuery string

// PaymentPatch godoc
// @Summary Patch payment
// @Schemes http
// @Description Patch payment by payment_id
// @Tags payment
// @Param id path int true "Payment ID"
// @Param client_id formData int false "Client ID"
// @Param room_id formData int false "Room ID"
// @Param payment_date formData string false "Date 'yyyy-mm-dd hh:mm:ss'"
// @Param payment_amount formData number false "Payment amount"
// @Produce json
// @Success 200 {object} types.APIResponse "Updated"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /payment/id/{id} [patch]
func RoutePaymentPatch(g *gin.Context) {
	var (
		apierr     *api_errors.APIError
		payment_id int64
		temp       string
		cache      = map[string]any{}
	)

	payment_id, apierr = validators.Int64("id", g.Param("id"), false)
	if apierr != nil {
		goto skip
	}

	temp = g.PostForm("client_id")
	if temp != "" {
		var v int64

		v, apierr = validators.Int64("client_id", temp, false)
		if apierr != nil {
			goto skip
		}

		cache["client_id"] = v
	}

	temp = g.PostForm("room_id")
	if temp != "" {
		var v int64

		v, apierr = validators.Int64("room_id", temp, false)
		if apierr != nil {
			goto skip
		}

		cache["room_id"] = v
	}

	temp = g.PostForm("payment_date")
	if temp != "" {
		var v time.Time

		v, apierr = validators.Date("payment_date", temp, false)
		if apierr != nil {
			goto skip
		}

		cache["payment_date"] = v.Format(validators.DATE_FORMAT)
	}

	temp = g.PostForm("payment_amount")
	if temp != "" {
		var v float64

		v, apierr = validators.Float64("payment_amount", temp, false)
		if apierr != nil {
			goto skip
		}

		cache["payment_amount"] = v
	}

skip:
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: apierr})
		return
	}

	var p types.Payment
	code, apierr := queryRow(&p, paymentScanRow, SQLPaymentGetByPaymentIDQuery, payment_id)
	if apierr != nil {
		g.JSON(code, types.APIResponse{Error: apierr})
		return
	}

	if _, ok := cache["client_id"]; !ok {
		cache["client_id"] = p.ClientID
	}
	if _, ok := cache["room_id"]; !ok {
		cache["room_id"] = p.RoomID
	}
	if _, ok := cache["payment_date"]; !ok {
		cache["payment_date"] = p.Date.Format(validators.DATE_FORMAT)
	}
	if _, ok := cache["payment_amount"]; !ok {
		cache["payment_amount"] = p.Amount
	}

	_, err := db.Exec(
		SQLPaymentPatchQuery,
		cache["client_id"],
		cache["room_id"],
		cache["payment_date"],
		cache["payment_amount"],
		payment_id,
	)
	if err != nil {
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	g.JSON(http.StatusOK, types.APIResponse{Message: "ok"})
}

func init() {
	r := api.Group("/payment")

	r.GET("/all", RoutePaymentGetAll)
	r.GET("/client/id/:id", RoutePaymentGetAllByClientID)
	r.GET("/room/id/:id", RoutePaymentGetAllByRoomID)
	r.GET("/date/:date", RoutePaymentGetByDate)
	r.POST("/date/range", RoutePaymentGetByDateRange)
	r.GET("/id/:id", RoutePaymentGetByID)
	r.DELETE("/id/:id", RoutePaymentDelete)
	r.PATCH("/id/:id", RoutePaymentPatch)
	r.POST("/new", RoutePaymentPostCreate)
}
