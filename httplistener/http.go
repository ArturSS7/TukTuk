package httplistener

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/telegrambot"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"html"
	"io/ioutil"
	"log"
	"time"
)

//start http server
func StartHTTP(db *sql.DB) {
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	e.GET("*", handleHTTP)
	e.POST("*", handleHTTP)
	e.PUT("*", handleHTTP)
	e.DELETE("*", handleHTTP)
	e.OPTIONS("*", handleHTTP)
	e.PATCH("*", handleHTTP)
	e.TRACE("*", handleHTTP)
	e.CONNECT("*", handleHTTP)
	e.Debug = true
	e.Logger.Fatal(e.Start(":80"))
}

//handler for http requests
//we make request to look as it should and store it in the database
//unfortunately i found no way to dump raw request with echo framework
func handleHTTP(c echo.Context) error {
	cc := c.(*database.DBContext)
	var result bool
	rows, err := cc.Db.Query("select exists (select id from dns_domains where domain = $1)", fmt.Sprintf("%s.", c.Request().Host))
	if err != nil {
		log.Println(err)
		return c.String(200, backend.RandStringBytes(8))
	}
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			log.Println(err)
			return c.String(200, backend.RandStringBytes(8))
		}
	}
	if result {
		request := new(bytes.Buffer)
		fmt.Fprintf(request, "%s %s %s\n", c.Request().Method, c.Request().URL, c.Request().Proto)
		fmt.Fprintf(request, "Host: %s\n", c.Request().Host)
		for i, v := range c.Request().Header {
			fmt.Fprintf(request, "%s: %s\n", i, v[0])
		}
		if c.Request().Body != nil {
			var bodyBytes []byte
			bodyBytes, err := ioutil.ReadAll(c.Request().Body)
			if err != nil {
				log.Println(err)
			}
			fmt.Fprintf(request, "\n%s", bodyBytes)
		}
		var lastInsertId int64 = 0
		err = cc.Db.QueryRow("insert into http (data, source_ip, time) values ($1, $2, $3) RETURNING id", html.EscapeString(request.String()), c.Request().RemoteAddr, time.Now().String()).Scan(&lastInsertId)

		if err != nil {
			log.Println(err)
		}

		//Send Alert to telegram
		telegrambot.BotSendAlert(html.EscapeString(request.String()), c.Request().RemoteAddr, time.Now().String(), "HTTP", lastInsertId)
	}
	return c.String(200, backend.RandStringBytes(8))
}
