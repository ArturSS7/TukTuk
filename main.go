package main

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/ftplistener"
	"TukTuk/httplistener"
)

func main() {
	//коннектим бдху
	db := database.Connect()
	//старт нттп для отстуков
	go httplistener.StartHTTP(db)
	//старт фтп для отстуков
	go ftplistener.StartFTP(db)
	//страт бека
	backend.StartBack(db)
}
