package backend

import (
	"TukTuk/database"
	"TukTuk/ftplistener"
	"TukTuk/plaintcplistener"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

type Request struct {
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

type TcpErr struct {
	Error error `json:"error"`
}

type TcpResult struct {
	Port    string `json:"port"`
	Success bool   `json:"success"`
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type user struct {
	username string
	password string
}

var credentials user

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
	secret := []byte("#JVb0VYu*3j!8oQmOtZK")
	e.Use(session.Middleware(sessions.NewCookieStore(secret)))
	e.Renderer = t
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	credentials.username = "dsec"
	credentials.password = "tuktuk"
	e.File("/", "frontend/index.html", loginRequired)
	e.File("/dns", "frontend/dns.html", loginRequired)
	e.Static("/static", "frontend/static/")
	e.GET("/api/:proto", getRequests, loginRequired)
	e.GET("/api/request/:proto", getRequest, loginRequired)
	e.GET("/api/dns/new", generateDomain, loginRequired)
	e.POST("/api/tcp/new", startPlainTCP, loginRequired)
	e.GET("/api/tcp/data", getTCPResults, loginRequired)
	e.POST("/api/tcp/shutdown", stopTCPServer, loginRequired)
	e.POST("/api/ftp/start", startFTP, loginRequired)
	e.POST("/api/ftp/shutdown", shutdownFTP, loginRequired)
	e.GET("/login", loginPage)
	e.POST("/login", handleLogin)
	e.GET("/api/dns/available", getAvailableDomains, loginRequired)
	e.Debug = true
	e.Logger.Fatal(e.Start(":1234"))
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
	default:
		return c.String(404, "Not Found")
	}
	limit := c.FormValue("limit")
	cc := c.(*database.DBContext)
	rows, err := cc.Db.Query("select data, source_ip, time from "+table+" order by id desc limit $1", limit)
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	rr := make([]Request, 0)
	for rows.Next() {
		r := Request{}
		err := rows.Scan(&r.Data, &r.SourceIp, &r.Time)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
		rr = append(rr, r)
	}
	return c.JSON(200, rr)
}

type Domain struct {
	Data string `json:"domain"`
}

var domain string

func generateDomain(c echo.Context) error {
	d := &Domain{}
	d.Data = RandStringBytes(8) + "." + domain
	cc := c.(*database.DBContext)
	_, err := cc.Db.Exec("insert into dns_domains (domain) values ($1)", d.Data)
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	return c.JSON(200, d)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandStringBytes(n int) string {
	rand.Seed(time.Now().Unix())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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

func loginPage(c echo.Context) error {
	if login := getLoginFromSession(c); login != "" {
		return c.Redirect(http.StatusFound, "/")
	} else {
		return c.Render(http.StatusOK, "login.html", nil)
	}
}

func handleLogin(c echo.Context) error {
	login := c.FormValue("username")
	password := c.FormValue("password")
	if login == credentials.username && password == credentials.password {
		sess := loginSession(c, login)
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return c.Render(http.StatusUnprocessableEntity, "login.html", "error")
		}
		return c.Redirect(http.StatusFound, "/")
	} else {
		return c.Render(http.StatusOK, "login.html",
			ErrorContext{"Incorrect username or password"},
		)
	}
}

func loginSession(c echo.Context, login string) *sessions.Session {
	sess, _ := session.Get("session", c)
	sess.Values["username"] = login
	sess.Options = &sessions.Options{
		Path: "/",
	}
	sess.Values["username"] = login
	return sess
}

func loginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if login := getLoginFromSession(c); login == "" {
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}

func getLoginFromSession(c echo.Context) string {
	sess, _ := session.Get("session", c)
	login, exists := sess.Values["username"]
	if !exists {
		return ""
	}
	return login.(string)
}

func getAvailableDomains(c echo.Context) error {
	cc := c.(*database.DBContext)
	rows, err := cc.Db.Query("select domain from dns_domains order by id desc")
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	dd := make([]Domain, 0)
	for rows.Next() {
		d := Domain{}
		err = rows.Scan(&d.Data)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
		dd = append(dd, d)
	}
	return c.JSON(200, dd)
}

func startPlainTCP(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		log.Println(err)
		return c.JSON(200, Result{Error: "true"})
	}
	if m["port"] == nil || m["port"] == "" {
		return c.JSON(200, Result{Error: "true"})
	}
	port := fmt.Sprintf("%v", m["port"])
	message := fmt.Sprintf("%v", m["message"])
	cc := c.(*database.DBContext)
	e := make(chan error)
	go func(r chan error) {
		err := plaintcplistener.StartTCP(cc.Db, message, port)
		if err != nil {
			r <- err
		}
	}(e)
	time.Sleep(2 * time.Second)
	select {
	case err := <-e:
		log.Println(err)
		close(e)
		er := &TcpErr{Error: err}
		return c.JSON(200, er)
	default:
		res := &TcpResult{Port: port, Success: true}
		close(e)
		return c.JSON(200, res)
	}
}

type TCPData struct {
	Data     string `json:"data"`
	SourceIP string `json:"source_ip"`
	Time     string `json:"time"`
	Port     string `json:"port"`
	Id       int    `json:"id"`
}

func getTCPResults(c echo.Context) error {
	port := c.QueryParam("port")
	cc := c.(*database.DBContext)
	var rows *sql.Rows
	var err error
	if port == "" {
		rows, err = cc.Db.Query("select * from tcp order by id desc")
	} else {
		rows, err = cc.Db.Query("select * from tcp where port = $1 order by id desc", port)
	}
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	tr := make([]TCPData, 0)
	for rows.Next() {
		t := TCPData{}
		err = rows.Scan(&t.Id, &t.Data, &t.SourceIP, &t.Time, &t.Port)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
		tr = append(tr, t)
	}
	return c.JSON(200, tr)
}

func stopTCPServer(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		log.Println(err)
		return c.JSON(200, Result{Error: "true"})
	}
	if m["port"] == nil || m["port"] == "" {
		return c.JSON(200, Result{Error: "true"})
	}
	port := fmt.Sprintf("%v", m["port"])
	if plaintcplistener.TCPServers[port] != nil {
		plaintcplistener.TCPServers[port].Stop()
		delete(plaintcplistener.TCPServers, port)
		return c.JSON(200, TcpResult{
			Port:    port,
			Success: true,
		})
	} else {
		return c.JSON(200, TcpResult{
			Port:    port,
			Success: false,
		})
	}
}

func startFTP(c echo.Context) error {
	cc := c.(*database.DBContext)
	e := make(chan error)
	go func(r chan error) {
		err := ftplistener.StartFTP(cc.Db)
		if err != nil {
			r <- err
		}
	}(e)
	time.Sleep(2 * time.Second)
	select {
	case err := <-e:
		log.Println(err)
		close(e)
		er := &TcpErr{Error: err}
		return c.JSON(200, er)
	default:
		res := &TcpResult{Port: "21", Success: true}
		close(e)
		return c.JSON(200, res)
	}
}

func shutdownFTP(c echo.Context) error {
	if ftplistener.FTPServer != nil {
		ftplistener.FTPServer.Stop()
		ftplistener.FTPServer = nil
		return c.JSON(200, TcpResult{
			Port:    "21",
			Success: true,
		})
	}
	return c.JSON(200, TcpResult{
		Port:    "21",
		Success: false,
	})
}
