package main

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/ftplistener"
	"TukTuk/httplistener"
	"TukTuk/httpslistener"
)

func main() {
	//connect to database
	db := database.Connect()
	//start http server
	go httplistener.StartHTTP(db)
	//start https server
	go httpslistener.StartHTTPS(db)
	//start ftp server
	go ftplistener.StartFTP(db)
	//start backend
	backend.StartBack(db)
}
