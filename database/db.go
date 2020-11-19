package database

import (
	"TukTuk/config"
	"database/sql"
	"fmt"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type DBContext struct {
	echo.Context
	Db *sql.DB
}

var DNSDB *sql.DB

func Connect() *sql.DB {
	ConnectString := "postgres://" + config.Settings.DBCredentials.Name + ":" + config.Settings.DBCredentials.Password + "@localhost/tuktuk?sslmode=disable"
	db, err := sql.Open("postgres", ConnectString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")
	DNSDB = db
	return db
}
