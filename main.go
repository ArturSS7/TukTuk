package main

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/dnslistener"
	"TukTuk/ftplistener"
	"TukTuk/httplistener"
	"TukTuk/smtplistener"
	"TukTuk/telegrambot"
)

func main() {
	//connect to database
	db := database.Connect()
	domain := "tt.pwn.bar."
	//start telegram bot
	telegrambot.BotStart()

	//start http server
	go httplistener.StartHTTP(db)

	//start https server
	//go httpslistener.StartHTTPS(db)

	//start ftp server
	go ftplistener.StartFTP(db)

	//start dns server
	go dnslistener.StartDNS(domain)

	//start smtp server
	go smtplistener.StartSMTP(db, domain)

	//start backend
	backend.StartBack(db, domain)

}
