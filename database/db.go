package database

import (
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
	db, err := sql.Open("postgres", "postgres://tuk:ZuppaSecurePwd@localhost/tuktuk?sslmode=disable")
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
