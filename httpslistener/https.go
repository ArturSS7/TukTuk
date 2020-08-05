package httpslistener

import (
	"TukTuk/database"
	"database/sql"
	"github.com/labstack/echo"
	"golang.org/x/crypto/acme/autocert"
)

func StartHTTPS(db *sql.DB) {
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.GET("*", handleHTTP)
	e.POST("*", handleHTTP)
	e.PUT("*", handleHTTP)
	e.DELETE("*", handleHTTP)
	e.OPTIONS("*", handleHTTP)
	e.PATCH("*", handleHTTP)
	e.TRACE("*", handleHTTP)
	e.CONNECT("*", handleHTTP)
	e.Debug = true
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}

func handleHTTP(c echo.Context) error {
	return c.String(200, "Q")
}
