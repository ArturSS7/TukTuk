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
	Id         int    `json:"id"`
	Data       string `json:"domain"`
	DeleteTime int64  `json:"delete_time"`
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
		_, err := cc.Db.Exec("insert into dns_domains (domain, delete_time) values ($1, 0)", d.Data)
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
		return c.JSON(200, d)
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

type DomainT struct {
	Id         int    `json:"id"`
	Data       string `json:"domain"`
	DeleteTime string `json:"delete_time"`
}

func getAvailableDomains(c echo.Context) error {
	cc := c.(*database.DBContext)
	_, err := cc.Db.Query("delete from dns_domains where delete_time < $1 and delete_time != 0", time.Now().Unix())
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	rows, err := cc.Db.Query("select id, domain, delete_time from dns_domains order by id desc")
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	dd := make([]Domain, 0)
	for rows.Next() {
		d := Domain{}
		err = rows.Scan(&d.Id, &d.Data, &d.DeleteTime)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
		dd = append(dd, d)
	}
	dt := make([]DomainT, 0)
	for _, v := range dd {
		d := DomainT{}
		if v.DeleteTime == 0 {
			d.DeleteTime = "Never"
		} else {
			d.DeleteTime = time.Unix(v.DeleteTime, 0).String()
		}
		d.Id = v.Id
		d.Data = v.Data
		dt = append(dt, d)
	}
	return c.JSON(200, dt)
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
