package main

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/dnslistener"
	"TukTuk/ftplistener"
	"TukTuk/httplistener"
	"TukTuk/httpslistener"
	"TukTuk/telegrambot"
)

func main() {
	//connect to database
	db := database.Connect()

	//start telegram bot
	telegrambot.BotStart("1351199153:AAEe1x20XTVb1Y4WWyp8DMzfOwcTca6rXE8", 367979213)

	//start http server
	go httplistener.StartHTTP(db)

	//start https server
	go httpslistener.StartHTTPS(db)

	//start ftp server
	go ftplistener.StartFTP(db)

	//start dns server
	go dnslistener.StartDNS(db)

	//start backend
	backend.StartBack(db)

}
