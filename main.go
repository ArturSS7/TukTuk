package main

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/dnslistener"
	"TukTuk/emailalert"
	"TukTuk/ftplistener"
	"TukTuk/httplistener"
	"TukTuk/httpslistener"
	"TukTuk/smtplistener"
	"TukTuk/telegrambot"
	"log"
)

func main() {
	//connect to database
	db := database.Connect()
	domain := "tt.pwn.bar."

	emailalert.Enabled = true
	telegrambot.Enabled = true
	if telegrambot.Enabled {
		//start telegram bot
		telegrambot.BotStart()
	}
	if err, res := emailalert.CheckConfig(); res && emailalert.Enabled {
		emailalert.GetClientToken()
	} else {
		log.Println(err)
		emailalert.Enabled = false
	}

	//start http server
	go httplistener.StartHTTP(db)

	//start https server
	go httpslistener.StartHTTPS(db)

	//start ftp server
	go ftplistener.StartFTP(db)

	//start dns server
	go dnslistener.StartDNS(domain)

	//start smtp server
	go smtplistener.StartSMTP(db, domain)

	//start backend
	backend.StartBack(db, domain)

}
