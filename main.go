package main

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/dnslistener"
	"TukTuk/ftplistener"
	"TukTuk/httplistener"
	"TukTuk/httpslistener"
	"TukTuk/smtplistener"
	"TukTuk/telegrambot"
)

func main() {
	//connect to database
	db := database.Connect()

	//start telegram bot
	telegrambot.BotStart()

	//start http server
	go httplistener.StartHTTP(db)

	//start https server
	go httpslistener.StartHTTPS(db)

	//start ftp server
	go ftplistener.StartFTP(db)

	//start dns server
	go dnslistener.StartDNS("tt.pwn.bar.")

	//start smtp server
	go smtplistener.StartSMTP(db, "localhost")

	//start backend
	backend.StartBack(db)

}
