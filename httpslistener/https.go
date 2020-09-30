package httpslistener

import (
	"TukTuk/backend"
	"TukTuk/config"
	"TukTuk/database"
	"TukTuk/discordbot"
	"TukTuk/emailalert"
	"TukTuk/telegrambot"
	"bytes"
	"database/sql"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"regexp"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/acme/autocert"
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
	e.HideBanner = true
	e.Debug = true
	e.Logger.Fatal(e.StartTLS(":443", config.Settings.HttpsCertPath.CertFile, config.Settings.HttpsCertPath.KeyFile))
}

func handleHTTPS(c echo.Context) error {
	cc := c.(*database.DBContext)
	var result bool
	domain := config.Settings.DomainConfig.Name[:len(config.Settings.DomainConfig.Name)-1]
	re := regexp.MustCompile(`([a-z0-9\-]+\.` + domain + `)`)
	d := re.Find([]byte(c.Request().Host))
	rows, err := cc.Db.Query("select exists (select id from dns_domains where domain = $1)", string(d)+".")
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

		//Send alert to Telegram
		telegrambot.BotSendAlert(html.EscapeString(request.String()), c.Request().RemoteAddr, time.Now().String(), "HTTPS", lastInsertId)
		//Send alert to email
		emailalert.SendEmailAlert("HTTPS Alert", "Remoute Address: "+c.Request().RemoteAddr+"\n+"+html.EscapeString(request.String())+"\n"+time.Now().String())
		//Send alert to Disord
		discordbot.BotSendAlert(html.EscapeString(request.String()), c.Request().RemoteAddr, time.Now().String(), "HTTPS", lastInsertId)

	}
	return c.String(200, backend.RandStringBytes(8))
}
