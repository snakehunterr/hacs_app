package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	api_errors "github.com/snakehunterr/hacs_app/db_types/errors"
)

var (
	logger          = log.New(os.Stdout, "[DB]", 0)
	logInfoMessage  = lipgloss.NewStyle().Background(lipgloss.Color("2")).Render("[INFO]")
	logErrorMessage = lipgloss.NewStyle().Background(lipgloss.Color("1")).Render("[ERROR]")
)

func logInfo(v ...any) {
	logger.Println(fmt.Sprintf(" %s %s: ", time.Now().Format("2006/01/02 15:04:05"), logInfoMessage) + fmt.Sprint(v...))
}

func logError(v ...any) {
	logger.Println(fmt.Sprintf(" %s %s: ", time.Now().Format("2006/01/02 15:04:05"), logErrorMessage) + fmt.Sprint(v...))
}

func queryRow[T any](
	dst T,
	fn func(T, *sql.Row) error,
	query string,
	a ...any,
) (code int, e *api_errors.APIError) {
	row := db.QueryRow(query, a...)

	if err := fn(dst, row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return http.StatusNotFound, api_errors.NewErrSQLNoRows("No rows")
		}

		logError("queryRow scanRowFunc() err:", err)
		return http.StatusInternalServerError, api_errors.NewErrSQLInternalError(err.Error())
	}

	return http.StatusOK, nil
}

func queryRows[T any](
	dst T,
	fn func(T, *sql.Rows) error,
	query string,
	a ...any,
) (code int, e *api_errors.APIError) {
	var (
		rows *sql.Rows
		err  error
	)

	if a == nil {
		rows, err = db.Query(query)
	} else {
		rows, err = db.Query(query, a...)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return http.StatusNotFound, api_errors.NewErrSQLNoRows("No rows")
		}

		logError("queryRows db.Query() err:", err)
		return http.StatusInternalServerError, api_errors.NewErrSQLInternalError(err.Error())
	}
	defer rows.Close()

	if err := fn(dst, rows); err != nil {
		logError("queryRows scanRowsFunc() err:", err)
		return http.StatusInternalServerError, api_errors.NewErrSQLInternalError(err.Error())
	}

	return http.StatusOK, nil
}
