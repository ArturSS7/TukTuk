package smblistener

import (
	"TukTuk/database"
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"time"
)

func StartSMBAccept(db *sql.DB) {
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	e.POST("/", acceptSMB)
	e.HideBanner = true
	e.Logger.Fatal(e.Start("127.0.0.1:5555"))
}

func acceptSMB(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		log.Println(err)
		return c.NoContent(500)
	}
	if m["data"] == "" || m["source_ip"] == "" {
		return c.NoContent(500)
	}
	cc := c.(*database.DBContext)
	_, err := cc.Db.Exec("insert into smb(data, source_ip, time) values ($1, $2, $3)", fmt.Sprintf("%v", m["data"]), fmt.Sprintf("%v", m["source_ip"]), time.Now().String())
	if err != nil {
		log.Println(err)
		return c.NoContent(500)
	}
	return c.NoContent(200)
}
