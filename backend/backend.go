package backend

import (
	"TukTuk/config"
	"TukTuk/database"
	"database/sql"
	"golang.org/x/crypto/acme/autocert"
	"html/template"
	"io"
	"log"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

type Request struct {
	Id       int    `json:"id"`
	Data     string `json:"data"`
	SourceIp string `json:"source_ip"`
	Time     string `json:"time"`
}

type SingleRequest struct {
	R     *Request
	Table string
}
type Result struct {
	Error string `json:"error"`
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type ErrorContext struct {
	Error string
}

//start backend
func StartBack(db *sql.DB, Domain string) {
	domain = Domain
	e := echo.New()
	//pass db pointer to echo handler
	t := &Template{
		templates: template.Must(template.ParseGlob("frontend/templates/*")),
	}
	secret := []byte(RandStringBytes(20))
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.Use(session.Middleware(sessions.NewCookieStore(secret)))
	e.Renderer = t
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	credentials.username = config.Settings.AdminCredentials.Username
	credentials.password = config.Settings.AdminCredentials.Password
	//e.Pre(middleware.HTTPSRedirect())
	e.File("/", "frontend/index.html", loginRequired)
	e.File("/tcp", "frontend/tcp.html", loginRequired)
	e.File("/dns", "frontend/dns.html", loginRequired)
	e.Static("/static", "frontend/static/")
	e.GET("/api/:proto", getRequests, loginRequired)
	e.GET("/api/request/:proto", getRequest, loginRequired)
	e.POST("/api/dns/new", generateDomain, loginRequired)
	e.POST("/api/dns/delete", deleteDomain, loginRequired)
	e.POST("/api/tcp/new", startPlainTCP, loginRequired)
	e.GET("/api/tcp/data", getTCPResults, loginRequired)
	e.POST("/api/tcp/shutdown", stopTCPServer, loginRequired)
	e.GET("/api/tcp/running", getRunningTCPServers, loginRequired)
	e.POST("/api/ftp/start", startFTP, loginRequired)
	e.POST("/api/ftp/shutdown", shutdownFTP, loginRequired)
	e.GET("/login", loginPage)
	e.POST("/login", handleLogin)
	e.GET("/api/dns/available", getAvailableDomains, loginRequired)
	e.POST("/api/smb/start", startSMBServer, loginRequired)
	e.POST("/api/smb/shutdown", stopSMBServer, loginRequired)
	e.HideBanner = true
	e.Debug = true
	e.Logger.Fatal(e.StartAutoTLS(":1234"))
	//e.Logger.Fatal(e.Start(":1234"))
}

//handler for getting requests from database

func getRequests(c echo.Context) error {
	table := ""
	switch c.Param("proto") {
	case "http":
		table = "http"
	case "ftp":
		table = "ftp"
	case "https":
		table = "https"
	case "dns":
		table = "dns"
	case "smtp":
		table = "smtp"
	case "ldap":
		table = "ldap"
	case "smb":
		table = "smb"
	default:
		return c.String(404, "Not Found")
	}
	limit := c.FormValue("limit")
	cc := c.(*database.DBContext)
	rows, err := cc.Db.Query("select id, data, source_ip, time from "+table+" order by id desc limit $1", limit)
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	rr := make([]Request, 0)
	for rows.Next() {
		r := Request{}
		err := rows.Scan(&r.Id, &r.Data, &r.SourceIp, &r.Time)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
		rr = append(rr, r)
	}
	return c.JSON(200, rr)
}

func getRequest(c echo.Context) error {
	table := ""
	switch c.Param("proto") {
	case "http":
		table = "http"
	case "ftp":
		table = "ftp"
	case "https":
		table = "https"
	case "dns":
		table = "dns"
	case "smtp":
		table = "smtp"
	case "ldap":
		table = "ldap"
	case "smb":
		table = "smb"
	default:
		return c.String(404, "Not Found")
	}
	id := c.QueryParam("id")
	cc := c.(*database.DBContext)
	var res bool
	rows, err := cc.Db.Query("select exists(select id from "+table+" where id = $1)", id)
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	for rows.Next() {
		err = rows.Scan(&res)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
	}
	if res {
		rows, err = cc.Db.Query("select data, source_ip, time from "+table+" where id = $1", id)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
		r := &Request{}
		for rows.Next() {
			err := rows.Scan(&r.Data, &r.SourceIp, &r.Time)
			if err != nil {
				log.Println(err)
				er := &Result{Error: "true"}
				return c.JSON(200, er)
			}
		}
		s := &SingleRequest{
			R:     r,
			Table: strings.ToTitle(table),
		}
		return c.Render(200, "request.html", s)
	}
	er := &Result{Error: "true"}
	return c.JSON(200, er)
}
