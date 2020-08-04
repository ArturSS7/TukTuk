package main

import (
	"TukTuk/backend"
	"TukTuk/database"
	"TukTuk/httplistener"
)

func main() {
	//коннектим бдху
	db := database.Connect()
	//страт бека
	go backend.StartBack(db)
	//стартр нттп для отстуков
	httplistener.StartHTTP(db)

}
