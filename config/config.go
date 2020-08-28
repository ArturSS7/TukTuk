//Domain config
//Telegram config
//GmailAPI config
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var Settings InitStruct

//StartInit domain and alert from Config.json.example
func StartInit() {
	parseConfig()
}

var ConfigPath string

func readConfig() []byte {
	var fileData []byte
	file, err := os.Open(ConfigPath)
	if err != nil {
		log.Fatalln(err)
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

	err := json.Unmarshal(b, &Settings)
	if err != nil {
		fmt.Println(err)
	}
}
