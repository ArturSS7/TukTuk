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
)

func main() {
	//connect to database
	db := database.Connect()
	domain := "tt.pwn.bar."

	//start telegram bot
	telegrambot.BotStart()

	emailalert.Enabled = true
	telegrambot.Enabled = true

	emailalert.EmailAlertStart()

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
