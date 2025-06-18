package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	_ "embed"

	"github.com/gin-gonic/gin"
	types "github.com/snakehunterr/hacs_db_types"
	api_errors "github.com/snakehunterr/hacs_db_types/errors"
	validators "github.com/snakehunterr/hacs_db_types/validators"
)

func roomScanRow(r *types.Room, row *sql.Row) error {
	return row.Scan(&r.ID, &r.ClientID, &r.PeopleCount, &r.Area, &r.LastEdited)
}

func roomScanRows(rs *[]types.Room, rows *sql.Rows) error {
	if rs == nil {
		return errors.New("*[]types.Room is nil")
	}

	_rs := *rs
	for rows.Next() {
		var r types.Room

		if err := rows.Scan(&r.ID, &r.ClientID, &r.PeopleCount, &r.Area, &r.LastEdited); err != nil {
			return err
		}

		_rs = append(_rs, r)
	}

	*rs = _rs
	return nil
}

//go:embed sql/room/room_get_all.sql
var SQLRoomGetAllQuery string

// RoomAll godoc
// @Summary Get all rooms
// @Schemes http
// @Description Get all rooms from MySQL
// @Tags room
// @Produce json
// @Success 200 {array} types.Room "ok"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /room/all [get]
func RouteRoomGetAll(g *gin.Context) {
	var rs []types.Room

	code, err := queryRows(&rs, roomScanRows, SQLRoomGetAllQuery)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(rs) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, rs)
}

//go:embed sql/room/room_get_by_id.sql
var SQLRoomGetByIDQuery string

// RoomByID godoc
// @Summary Get room by room_id
// @Schemes http
// @Description Get room by room_id
// @Tags room
// @Param id path int true "Room ID"
// @Produce json
// @Success 200 {object} types.Room "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /room/id/{id} [get]
func RouteRoomGetByID(g *gin.Context) {
	id := g.Param("id")
	if _, apierr := validators.Int64("id", id, false); apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: apierr})
		return
	}

	var r types.Room
	code, err := queryRow(&r, roomScanRow, SQLRoomGetByIDQuery, id)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	g.JSON(http.StatusOK, r)
}

//go:embed sql/room/room_get_by_client_id.sql
var SQLRoomGetByClientIDQuery string

// RoomByClientID godoc
// @Summary Get rooms by client_id
// @Schemes http
// @Description Get rooms by client_id
// @Tags room
// @Param id path int true "Client ID"
// @Produce json
// @Success 200 {array} types.Room "ok"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "No rows"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /room/client/id/{id} [get]
func RouteRoomGetByClientID(g *gin.Context) {
	id := g.Param("id")
	if _, err := validators.Int64("id", id, false); err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
		return
	}

	var rs []types.Room
	code, err := queryRows(&rs, roomScanRows, SQLRoomGetByClientIDQuery, id)
	if err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if len(rs) == 0 {
		g.JSON(http.StatusNotFound, types.APIResponse{
			Error: api_errors.NewErrSQLNoRows("No rows"),
		})
		return
	}

	g.JSON(http.StatusOK, rs)
}

//go:embed sql/room/room_insert.sql
var SQLRoomPostCreateQuery string

// RoomCreate godoc
// @Summary Create new room
// @Schemes http
// @Description Create new room
// @Tags room
// @Param id path int true "Room ID"
// @Param client_id formData int true "Client ID"
// @Param room_people_count formData int true "People living in room"
// @Param room_area formData number true "Room area"
// @Produce json
// @Success 201 {object} types.APIResponse "Created"
// @Failure 400 {object} types.APIResponse "Missing parameter"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /room/id/{id} [post]
func RouteRoomPostCreate(g *gin.Context) {
	var (
		apierr       *api_errors.APIError
		room_id      int64
		client_id    int64
		room_area    float64
		people_count uint8
	)

	room_id, apierr = validators.Int64("id", g.Param("id"), false)
	if apierr != nil {
		goto skip
	}

	client_id, apierr = validators.Int64("client_id", g.PostForm("client_id"), true)
	if apierr != nil {
		goto skip
	}

	room_area, apierr = validators.Float64("room_area", g.PostForm("room_area"), true)
	if apierr != nil {
		goto skip
	}

	people_count, apierr = validators.Uint8("room_people_count", g.PostForm("room_people_count"), true)

skip:
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{
			Error: apierr,
		})
		return
	}

	_, err := db.Exec(SQLRoomPostCreateQuery, room_id, client_id, people_count, room_area)
	if err != nil {
		logError("db.Exec(): ", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	logInfo(fmt.Sprintf("Created room: %#v", types.Room{ID: room_id, ClientID: client_id, Area: room_area, PeopleCount: people_count}))
	g.JSON(http.StatusCreated, types.APIResponse{Message: "ok"})
}

//go:embed sql/room/room_delete.sql
var SQLRoomDeleteQuery string

// RoomDelete godoc
// @Summary Delete room by room_id
// @Schemes http
// @Description Delete room by room_id
// @Tags room
// @Param id path int true "Room ID"
// @Produce json
// @Success 200 {object} types.APIResponse "Deleted"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /room/id/{id} [delete]
func RouteRoomDelete(g *gin.Context) {
	id := g.Param("id")
	if _, err := validators.Int64("id", id, false); err != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: err})
	}

	_, err := db.Exec(SQLRoomDeleteQuery, id)
	if err != nil {
		logError("db.Exec():", err)
		g.JSON(http.StatusInternalServerError, types.APIResponse{
			Error: api_errors.NewErrSQLInternalError(err.Error()),
		})
		return
	}

	logInfo("Deleted room record with room_id:", id)
	g.JSON(http.StatusOK, types.APIResponse{Message: "ok"})
}

//go:embed sql/room/room_patch.sql
var SQLRoomPatchQuery string

// RoomPatch godoc
// @Summary Patch room
// @Schemes http
// @Description Patch room by room_id
// @Tags room
// @Param id path int true "Room ID"
// @Param client_id formData int false "Client ID"
// @Param room_area formData number false "Room area"
// @Param room_people_count formData int false "People living in room"
// @Produce json
// @Success 200 {object} types.APIResponse "Updated"
// @Failure 400 {object} types.APIResponse "Incorrect parameter"
// @Failure 404 {object} types.APIResponse "Record not founded in DB"
// @Failure 500 {object} types.APIResponse "Internal server error"
// @Router /room/id/{id} [patch]
func RouteRoomPatch(g *gin.Context) {
	var (
		apierr  *api_errors.APIError
		room_id int64
		temp    string
		cache   = map[string]any{}
	)

	room_id, apierr = validators.Int64("id", g.Param("id"), false)
	if apierr != nil {
		goto skip
	}

	temp = g.PostForm("client_id")
	if temp != "" {
		var client_id int64

		client_id, apierr = validators.Int64("client_id", temp, false)
		if apierr != nil {
			goto skip
		}
		cache["client_id"] = client_id
	}

	temp = g.PostForm("room_area")
	if temp != "" {
		var room_area float64

		room_area, apierr = validators.Float64("room_area", temp, false)
		if apierr != nil {
			goto skip
		}
		cache["room_area"] = room_area
	}

	temp = g.PostForm("room_people_count")
	if temp != "" {
		var people_count uint8

		people_count, apierr = validators.Uint8("room_people_count", temp, false)
		if apierr != nil {
			goto skip
		}
		cache["people_count"] = people_count
	}

skip:
	if apierr != nil {
		g.JSON(http.StatusBadRequest, types.APIResponse{Error: apierr})
		return
	}

	var r types.Room

	if code, err := queryRow(&r, roomScanRow, SQLRoomGetByIDQuery, room_id); err != nil {
		g.JSON(code, types.APIResponse{Error: err})
		return
	}

	if _, ok := cache["client_id"]; !ok {
		cache["client_id"] = r.ClientID
	}
	if _, ok := cache["room_area"]; !ok {
		cache["room_area"] = r.Area
	}
	if _, ok := cache["people_count"]; !ok {
		cache["people_count"] = r.PeopleCount
	}

	_, err := db.Exec(
		SQLRoomPatchQuery,
		cache["client_id"],
		cache["people_count"],
		cache["room_area"],
		room_id,
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
	r := api.Group("/room")

	r.GET("/all", RouteRoomGetAll)
	r.GET("/id/:id", RouteRoomGetByID)
	r.POST("/id/:id", RouteRoomPostCreate)
	r.DELETE("/id/:id", RouteRoomDelete)
	r.PATCH("/id/:id", RouteRoomPatch)
	r.GET("/client/id/:id", RouteRoomGetByClientID)
}
