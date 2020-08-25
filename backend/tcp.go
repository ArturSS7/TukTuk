package backend

import (
	"TukTuk/database"
	"TukTuk/plaintcplistener"
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"math/rand"
	"time"
)

type TcpErr struct {
	Error error `json:"error"`
}

type TcpResult struct {
	Port    string `json:"port"`
	Success bool   `json:"success"`
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

type TcpServer struct {
	Port string `json:"port"`
}

func getRunningTCPServers(c echo.Context) error {
	ts := make([]TcpServer, 0)
	for v := range plaintcplistener.TCPServers {
		t := TcpServer{Port: v}
		ts = append(ts, t)
	}
	return c.JSON(200, ts)
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
