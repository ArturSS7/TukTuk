package backend

import (
	"TukTuk/database"
	"database/sql"
	"github.com/labstack/echo"
	"html/template"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"
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

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//start backend
func StartBack(db *sql.DB) {
	e := echo.New()
	//pass db pointer to echo handler
	t := &Template{
		templates: template.Must(template.ParseGlob("frontend/templates/*")),
	}
	e.Renderer = t
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	e.File("/", "frontend/index.html")
	e.Static("/static", "frontend/static/")
	e.GET("/api/:proto", getRequests)
	e.GET("/api/request/:proto", getRequest)
	e.GET("/api/dns/new", generateDomain)
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

func generateDomain(c echo.Context) error {
	d := &Domain{}
	d.Data = RandStringBytes(8) + ".tt.pwn.bar"
	cc := c.(*database.DBContext)
	_, err := cc.Db.Exec("insert into dns_domains (domain) values ($1)", d.Data+".")
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
