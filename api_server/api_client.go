package main

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	types "github.com/snakehunterr/hacs_app/db_types"
	api_errors "github.com/snakehunterr/hacs_app/db_types/errors"
	validators "github.com/snakehunterr/hacs_app/db_types/validators"
)

func clientScanRow(c *types.Client, row *sql.Row) error {
	return row.Scan(&c.ID, &c.Name, &c.IsAdmin, &c.LastEdited)
}

func clientScanRows(cs *[]types.Client, rows *sql.Rows) error {
	if cs == nil {
		return errors.New("*[]types.Client is nil")
	}
	_cs := *cs
	for rows.Next() {
		var c types.Client

		if err := rows.Scan(&c.ID, &c.Name, &c.IsAdmin, &c.LastEdited); err != nil {
			return err
		}

		_cs = append(_cs, c)
	}
	*cs = _cs
	return nil
}

//go:embed sql/client/client_get_all.sql
var SQLClientGetAllQuery string

// ClientAll godoc
// @Summary Get all clients
// @Schemes http
// @Description Get all clients
// @Tags client
// @Produce json
// @Success 200 {array} types.Client "ok"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /client/all [get]
func RouteClientGetAll(g *gin.Context) {
	var cs []types.Client

	code, err := queryRows(&cs, clientScanRows, SQLClientGetAllQuery)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(cs) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(code, cs)
}

//go:embed sql/client/client_get_admins.sql
var SQLClientGetAdminsQuery string

// ClientAllAdmins godoc
// @Summary Get all admin clients
// @Schemes http
// @Description Get all admin clients
// @Tags client
// @Produce json
// @Success 200 {array} types.Client "ok"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /client/admins [get]
func RouteClientGetAdmins(g *gin.Context) {
	var cs []types.Client

	code, err := queryRows(&cs, clientScanRows, SQLClientGetAdminsQuery)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(cs) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(code, cs)
}

//go:embed sql/client/client_get_by_id.sql
var SQLClientGetByIDQuery string

// ClientByTelegramID godoc
// @Summary Get client by telegram ID
// @Schemes http
// @Description Get client by telegram ID
// @Tags client
// @Param id path int true "Telegram ID"
// @Produce json
// @Success 200 {object} types.Client "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /client/id/{id} [get]
func RouteClientGetByID(g *gin.Context) {
	id := g.Param("id")
	if _, err := validators.Int64("id", id, false); err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	var c types.Client
	code, err := queryRow(&c, clientScanRow, SQLClientGetByIDQuery, id)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	g.JSON(code, c)
}

//go:embed sql/client/client_get_by_name.sql
var SQLClientGetByNameQuery string

// ClientByName godoc
// @Summary Get clients by client_name
// @Schemes http
// @Description Get clients by client_name
// @Tags client
// @Param name path string true "Client name"
// @Produce json
// @Success 200 {array} types.Client "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /client/name/{name} [get]
func RouteClientGetByName(g *gin.Context) {
	var cs []types.Client

	code, err := queryRows(&cs, clientScanRows, SQLClientGetByNameQuery, g.Param("name"))
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(cs) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(code, cs)
}

//go:embed sql/client/client_insert.sql
var SQLClientPostCreateQuery string

// ClientCreate godoc
// @Summary Create new client
// @Schemes http
// @Description Create new client
// @Tags client
// @Param id path int true "Client telegram ID"
// @Param client_name formData string true "Client name"
// @Param is_admin formData bool true "Is admin"
// @Produce json
// @Success 201 {object} types.Client "New client"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /client/id/{id} [post]
func RouteClientPostCreate(g *gin.Context) {
	var (
		apierr      *api_errors.APIError
		client_id   int64
		client_name string
		is_admin    bool
	)

	client_id, apierr = validators.Int64("id", g.Param("id"), false)
	if apierr != nil {
		goto skip
	}

	client_name = g.PostForm("client_name")
	if len(client_name) == 0 {
		apierr = api_errors.NewErrEmptyParam("client_name")
		goto skip
	}

	is_admin, apierr = validators.Bool("is_admin", g.PostForm("is_admin"), true)

skip:
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{
			Error: apierr,
		})
		return
	}

	_, err := db.Exec(SQLClientPostCreateQuery, client_id, client_name, is_admin)
	if err != nil {
		logError("db.Exec():", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	logInfo(fmt.Sprintf("Created new client: %#v", types.Client{
		ID:      client_id,
		Name:    client_name,
		IsAdmin: is_admin,
	}))
	g.JSON(http.StatusCreated, types.APIResponse{Message: "ok"})
}

//go:embed sql/client/client_delete.sql
var SQLClientDeleteQuery string

// ClientDelete godoc
// @Summary Delete client
// @Schemes http
// @Description Delete client by telegram ID
// @Tags client
// @Param id path int true "Client telegram ID"
// @Produce json
// @Success 200 {object} types.APIResponse "Deleted"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /client/id/{id} [delete]
func RouteClientDelete(g *gin.Context) {
	id, apierr := validators.Int64("id", g.Param("id"), false)
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: apierr})
		return
	}

	_, err := db.Exec(SQLClientDeleteQuery, id)
	if err != nil {
		logError("db.Exec():", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	logInfo("Deleted client with client_id:", id)
	g.JSON(http.StatusOK, types.APIResponse{Message: "ok"})
}

//go:embed sql/client/client_patch.sql
var SQLClientPatchQuery string

// ClientPatch godoc
// @Summary Patch client
// @Schemes http
// @Description Patch client by client_id
// @Tags client
// @Param id path int true "Client ID"
// @Param client_name formData string false "Client name"
// @Param is_admin formData bool false "is admin"
// @Produce json
// @Success 200 {object} types.APIResponse "Updated"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /client/id/{id} [patch]
func RouteClientPatch(g *gin.Context) {
	var (
		apierr    *api_errors.APIError
		client_id int64
		temp      string
		cache     = map[string]any{}
	)

	client_id, apierr = validators.Int64("id", g.Param("id"), false)
	if apierr != nil {
		goto skip
	}

	temp = g.PostForm("client_name")
	if temp != "" {
		cache["client_name"] = temp
	}

	temp = g.PostForm("is_admin")
	if temp != "" {
		var is_admin bool

		is_admin, apierr = validators.Bool("is_admin", temp, false)
		if apierr != nil {
			goto skip
		}
		cache["is_admin"] = is_admin
	}

skip:
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{
			Error: apierr,
		})
		return
	}

	var c types.Client
	code, apierr := queryRow(&c, clientScanRow, SQLClientGetByIDQuery, client_id)
	if apierr != nil {
		g.JSON(code, types.APIResponse{Error: apierr})
		return
	}

	if _, ok := cache["client_name"]; !ok {
		cache["client_name"] = c.Name
	}
	if _, ok := cache["is_admin"]; !ok {
		cache["is_admin"] = c.IsAdmin
	}

	_, err := db.Exec(SQLClientPatchQuery, cache["client_name"], cache["is_admin"], client_id)
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
	r := api.Group("/client")

	r.GET("/all", RouteClientGetAll)
	r.GET("/admins", RouteClientGetAdmins)
	r.GET("/name/:name", RouteClientGetByName)
	r.GET("/id/:id", RouteClientGetByID)
	r.POST("/id/:id", RouteClientPostCreate)
	r.DELETE("/id/:id", RouteClientDelete)
	r.PATCH("/id/:id", RouteClientPatch)
}
