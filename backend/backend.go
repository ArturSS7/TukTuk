package backend

import (
	"TukTuk/database"
	"database/sql"
	"github.com/labstack/echo"
	"log"
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
func StartBack(db *sql.DB) {
	e := echo.New()
	//pass db pointer to echo handler
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	e.GET("/api/:proto", getRequests)
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
