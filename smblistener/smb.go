package smblistener

import (
	"TukTuk/database"
	"TukTuk/emailalert"
	"TukTuk/telegrambot"
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/labstack/echo"
)

var CMD = exec.Command("python3", "smblistener/impacket_smb/smb.py")

func StartSMBAccept(db *sql.DB) {
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	e.POST("/", acceptSMB)
	e.HideBanner = true
	e.Logger.Fatal(e.Start("127.0.0.1:5555"))
}

func acceptSMB(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		log.Println(err)
		return c.NoContent(500)
	}
	if m["data"] == "" || m["source_ip"] == "" {
		return c.NoContent(500)
	}
	cc := c.(*database.DBContext)
	var lastInsertId int64 = 0
	err := cc.Db.QueryRow("insert into smb(data, source_ip, time) values ($1, $2, $3) RETURNING id", fmt.Sprintf("%v", m["data"]), fmt.Sprintf("%v", m["source_ip"]), time.Now().String()).Scan(&lastInsertId)
	//Send Alert to telegram
	telegrambot.BotSendAlert(fmt.Sprintf("%v", m["data"]), fmt.Sprintf("%v", m["source_ip"]), time.Now().String(), "SMB", lastInsertId)
	//Send Alert to email
	emailalert.SendEmailAlert("SMB Alert", fmt.Sprintf("%v", m["source_ip"])+"\n\n"+fmt.Sprintf("%v", m["data"]))
	if err != nil {
		log.Println(err)
		return c.NoContent(500)
	}
	return c.NoContent(200)
}
