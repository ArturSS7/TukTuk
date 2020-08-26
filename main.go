package main

import (
	"TukTuk/backend"
	"TukTuk/config"
	"TukTuk/database"
	"TukTuk/dnslistener"
	"TukTuk/emailalert"
	"TukTuk/httplistener"
	"TukTuk/httpslistener"
	"TukTuk/smblistener"
	"TukTuk/smtplistener"
	"TukTuk/telegrambot"
	"log"
	"os/exec"
)

func main() {
	config.StartInit()
	domain := config.Settings.DomainConfig.Name

	//start telegram bot
	telegrambot.BotStart()
	emailalert.EmailAlertStart(config.Settings.EmailAlert.Enabled, config.Settings.EmailAlert.To)

	//connect to database
	db := database.Connect()

	//start http server
	go httplistener.StartHTTP(db)

	//start https server
	go httpslistener.StartHTTPS(db)

	//start dns server
	go dnslistener.StartDNS(domain)

	//start smtp server
	go smtplistener.StartSMTP(db, domain)

	//start ldap server
	//	go ldaplistener.StartLDAP(domain)

	//start smb
	go smblistener.StartSMBAccept(db)
	cmd := exec.Command("python3", "smblistener/impacket_smb/smb.py")
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	//start backend
	backend.StartBack(db, domain)

}
