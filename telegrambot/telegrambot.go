package telegrambot

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

type content struct {
	data      string
	source_ip string
	time      string
}
type Option struct {
	token  string
	ChatID int64
}

var opt Option

func BotStart(_token string, _CharID int64) {
	opt = Option{_token, _CharID}
}

// 0 - Short, 1 - Long, Default = Short
func GetShort() {
	Setting = Short
}
func GetLong() {
	Setting = Long
}

const (
	Short = iota
	Long
)

var Setting int = Short

func BotSendAlert(data, source_ip, time, ProtocolName string) {
	bot, err := tgbotapi.NewBotAPI(opt.token)
	if err != nil {
		log.Panic(err)
	}
	_cont := content{data, source_ip, time}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//responce := tgbotapi.NewMessage(opt.ChatID, data+" "+source_ip+" "+time)
	responce := tgbotapi.NewMessage(opt.ChatID, MessageFormation(_cont, ProtocolName))
	bot.Send(responce)

}

//BotSendAlert_BD function start the bot and sends the message read from the database
func BotSendAlert_BD(tableName string, id int, db *sql.DB) {
	bot, err := tgbotapi.NewBotAPI(opt.token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	responce := tgbotapi.NewMessage(opt.ChatID, readDB(tableName, id, db))
	bot.Send(responce)

}

func readDB(tableName string, id int, db *sql.DB) string {
	str := "SELECT data, source_ip, time from " + tableName + " WHERE " + strconv.Itoa(id) + "= id"
	rows, err := db.Query(str)
	defer rows.Close()
	contents := []content{}
	if err != nil {
		panic(err)
	}
	rows.Next()
	p := content{}
	err = rows.Scan(&p.data, &p.source_ip, &p.time)
	if err != nil {
		fmt.Println(err)

	} else {
		contents = append(contents, p)
	}
	return MessageFormation(contents[0], tableName)
}

func MessageFormation(ContentFormation content, ProtocolName string) string {
	if Setting == Short {
		return "Catched " + ProtocolName + " request from IP: " + parsePort(ContentFormation.source_ip)
	} else {
		return ContentFormation.data + "\n" + parsePort(ContentFormation.source_ip) + "\n" + ContentFormation.time
	}
}

func parsePort(str string) string {
	re := regexp.MustCompile(":")
	return re.Split(str, -1)[0]
}
