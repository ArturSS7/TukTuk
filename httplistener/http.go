package httplistener

import (
	"TukTuk/database"
	"TukTuk/telegrambot"
	"bytes"
	"database/sql"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"time"

	"github.com/labstack/echo"
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
	cc := c.(*database.DBContext)
	_, err := cc.Db.Exec("insert into http (data, source_ip, time) values ($1, $2, $3)", html.EscapeString(request.String()), c.Request().RemoteAddr, time.Now().String())

	if err != nil {
		log.Println(err)
	}
	//Send Alert to telegram
	telegrambot.BotSendAlert("1351199153:AAEe1x20XTVb1Y4WWyp8DMzfOwcTca6rXE8", 367979213, html.EscapeString(request.String()), c.Request().RemoteAddr, time.Now().String())

	return c.String(200, "TukTuk callback server")
}
