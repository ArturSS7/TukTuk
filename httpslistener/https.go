package httpslistener

import (
	"TukTuk/database"
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
	_, err := cc.Db.Exec("insert into https (data, source_ip, time) values ($1, $2, $3)", html.EscapeString(request.String()), c.Request().RemoteAddr, time.Now().String())
	if err != nil {
		log.Println(err)
	}
	return c.String(200, "TukTuk callback server")
}
