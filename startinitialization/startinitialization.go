//Domen name
//Telegram config
//GmailAPI config
package startinitialization

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var Settings InitStruct

//StartInit domain and alert from Config.json
func StartInit() {
	parseConfig()
}

func readConfig() []byte {
	var fileData []byte
	file, err := os.Open("startinitialization/Config.json")
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
