package backend

import (
	"TukTuk/database"
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/labstack/echo"
)

type Request struct {
	Data     string `json:"data"`
	SourceIp string `json:"source_ip"`
	Time     string `json:"time"`
}

type Result struct {
	Error string `json:"error"`
}

//start backend
func StartBack(db *sql.DB, Domain string) {
	domain = Domain
	e := echo.New()
	//pass db pointer to echo handler
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	e.GET("/api/:proto", getRequests)
	e.GET("/api/dns/new", generateDomain)
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
	rows, err := cc.Db.Query("select data, source_ip, time from "+table+" order by id limit $1", limit)
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
