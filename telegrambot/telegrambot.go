package telegrambot

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
	Token       string
	ChatID      int64
	LenghtAlert string
}

//BotStart Start Telegram Bot
func BotStart() {
	parseConfig()
}

var SettingBot Option

func BotSendAlert(data, source_ip, time, ProtocolName string, id int64) {
	bot, err := tgbotapi.NewBotAPI(SettingBot.Token)
	if err != nil {
		log.Panic(err)
	}
	_cont := content{data, source_ip, time}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	responce := tgbotapi.NewMessage(SettingBot.ChatID, messageFormation(_cont, ProtocolName, id))
	bot.Send(responce)

}

//BotSendAlert_BD function start the bot and sends the message read from the database
func BotSendAlert_BD(tableName string, db *sql.DB, id int64) {
	bot, err := tgbotapi.NewBotAPI(SettingBot.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	responce := tgbotapi.NewMessage(SettingBot.ChatID, readDB(tableName, db, id))
	bot.Send(responce)

}

func readDB(tableName string, db *sql.DB, id int64) string {
	rows, err := db.Query("SELECT data, source_ip, time from $1 WHERE $2= id", tableName, id)
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
	return messageFormation(contents[0], tableName, id)
}

func messageFormation(ContentFormation content, ProtocolName string, id int64) string {

	if SettingBot.LenghtAlert == "Long" {
		return ContentFormation.data + "\n" + parsePort(ContentFormation.source_ip) + "\n" + ContentFormation.time + "\n\nLink: http://127.0.0.1:1234/api/request/http?id=" + strconv.Itoa(int(id))
	}
	return "Catched " + ProtocolName + " request from IP: " + parsePort(ContentFormation.source_ip) + "\n\nLink: http://pwn.bar:1234/api/request/http?id=" + strconv.Itoa(int(id))
}

func parsePort(str string) string {
	re := regexp.MustCompile(":")
	return re.Split(str, -1)[0]
}

func readConfig() []byte {
	var fileData []byte
	file, err := os.Open("telegrambot/Config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	data := make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}

		fileData = append(fileData, data[:n]...)

	}
	return fileData
}

func parseConfig() {
	b := readConfig()

	err := json.Unmarshal(b, &SettingBot)
	if err != nil {
		panic(err)
	}
}

func getIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsUnspecified() {
			if ipnet.IP.To4() != nil {
				os.Stdout.WriteString(ipnet.IP.String() + "\n")
			}
		}
	}
}
