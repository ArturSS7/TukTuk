package httplistener

import (
	"TukTuk/database"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"html"
	"io/ioutil"
	"log"
	"time"
)

//стратуем сервер для ловли нттп отстуков
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

//хендел для ловли ннтп отстуков, тут собираем нам реквест в человеческий вид и кладем в базу
//к сожалению сдампить в raw виде нельзя, но собирается тоже нормально
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
	return c.String(200, "TukTuk callback server")
}
