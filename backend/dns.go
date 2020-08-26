package backend

import (
	"TukTuk/database"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"strconv"
	"time"
)

type Domain struct {
	Data string `json:"domain"`
}

var domain string

var deleteTimes = []int64{3600, 86400, 604800, 2629743}

func generateDomain(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		log.Println(err)
		return c.JSON(200, Result{Error: "true"})
	}
	cc := c.(*database.DBContext)
	if m["delete_time"] == "" {
		d := &Domain{}
		d.Data = RandStringBytes(8) + "." + domain
		_, err := cc.Db.Exec("insert into dns_domains (domain) values ($1)", d.Data)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
		return c.JSON(200, d)
	}
	dT, err := strconv.ParseInt(fmt.Sprintf("%v", m["delete_time"]), 10, 64)
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	if Find(deleteTimes, dT) {
		d := &Domain{}
		d.Data = RandStringBytes(8) + "." + domain
		_, err := cc.Db.Exec("insert into dns_domains (domain, delete_time) values ($1, $2)", d.Data, time.Now().Unix()+dT)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
	}
	return c.JSON(200, Result{Error: "true"})
}

func Find(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func getAvailableDomains(c echo.Context) error {
	cc := c.(*database.DBContext)
	_, err := cc.Db.Query("delete from dns_domains where delete_time < $1 and delete_time is not null", time.Now().Unix())
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
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

func deleteDomain(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		log.Println(err)
		return c.JSON(200, Result{Error: "true"})
	}
	if m["domain"] == "" {
		return c.JSON(200, Result{Error: "true"})
	}
	cc := c.(*database.DBContext)
	_, err := cc.Db.Query("delete from dns_domains where domain = $1", fmt.Sprintf("%v", m["domain"]))
	if err != nil {
		log.Println(err)
		return c.JSON(200, Result{Error: "true"})
	}
	return c.JSON(200, Result{Error: "false"})
}
