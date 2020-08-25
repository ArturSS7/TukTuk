package main

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/dnslistener"
	"TukTuk/emailalert"
	"TukTuk/httplistener"
	"fmt"

	//"TukTuk/httpslistener"

	//"TukTuk/httpslistener"
	"TukTuk/config"
	"TukTuk/smtplistener"
	"TukTuk/telegrambot"
)

func main() {
	config.StartInit()
	domain := config.Settings.DomainConfig.Name
	fmt.Println(config.Settings.DomainConfig)
	//start telegram bot
	telegrambot.BotStart()
	emailalert.EmailAlertStart(config.Settings.EmailAlert.Enabled, config.Settings.EmailAlert.To)
	fmt.Println(config.Settings)
	//fmt.Println(config.Settings.DomainConfig.Name[:len(config.Settings.DomainConfig.Name)-1])
	//connect to database
	db := database.Connect()

	//start http server
	go httplistener.StartHTTP(db)

	//start https server
	//go httpslistener.StartHTTPS(db)

	//start dns server
	go dnslistener.StartDNS(domain)

	//start smtp server
	go smtplistener.StartSMTP(db, domain)

	//start backend
	backend.StartBack(db, domain)

}
