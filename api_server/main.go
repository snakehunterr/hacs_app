package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/snakehunterr/hacs_app/api_server/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title HACS database API
// @version 1.0
// @BasePath /api
// @produce json

var (
	db  *sql.DB
	e   *gin.Engine      = gin.Default()
	api *gin.RouterGroup = e.Group("/api")
)

func main() {
	db = openDB()
	api = e.Group("/api")

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", os.Getenv("DBAPI_SERVER_PORT"))

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	e.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "200 OK",
		})
	})

	if err := e.Run(fmt.Sprintf(":%s", os.Getenv("DBAPI_SERVER_PORT"))); err != nil {
		panic(err)
	}
}

func openDB() *sql.DB {
	for {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"),
		)
		logInfo("DSN: ", dsn)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			logError(err)
			time.Sleep(time.Second)
			continue
		}

		if err := db.Ping(); err != nil {
			logError(err)
			time.Sleep(time.Second * 3)
			continue
		}

		logInfo("Connected to:", dsn)

		return db
	}
}
