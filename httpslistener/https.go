package httpslistener

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/telegrambot"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/crypto/acme/autocert"
	"html"
	"io/ioutil"
	"log"
	"time"
)

//same as https, just added cert
func StartHTTPS(db *sql.DB) {
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.GET("*", handleHTTPS)
	e.POST("*", handleHTTPS)
	e.PUT("*", handleHTTPS)
	e.DELETE("*", handleHTTPS)
	e.OPTIONS("*", handleHTTPS)
	e.PATCH("*", handleHTTPS)
	e.TRACE("*", handleHTTPS)
	e.CONNECT("*", handleHTTPS)
	e.Debug = true
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}

func handleHTTPS(c echo.Context) error {
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
		cc := c.(*database.DBContext)
		var lastInsertId int64 = 0
		err := cc.Db.QueryRow("insert into https (data, source_ip, time) values ($1, $2, $3)  RETURNING id", html.EscapeString(request.String()), c.Request().RemoteAddr, time.Now().String()).Scan(&lastInsertId)
		if err != nil {
			log.Println(err)
		}

		//Send Alert to telegram
		telegrambot.BotSendAlert(html.EscapeString(request.String()), c.Request().RemoteAddr, time.Now().String(), "HTTPS", lastInsertId)
	}
	return c.String(200, backend.RandStringBytes(8))
}
