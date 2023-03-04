package database

import (
	"database/sql"

	"github.com/fundraising/rest-api/config"
	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB
var err error

func Initialize() {

	conf := config.GetConfig()

	connectionString := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME

	database, err = sql.Open("mysql", connectionString)

	if err != nil {
		panic(err.Error())
	}

	err = database.Ping()

	if err != nil {
		panic("DSN invalid")
	}
}

func CreateConnection() *sql.DB {
	return database
}
