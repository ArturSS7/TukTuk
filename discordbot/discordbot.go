package discordbot

import (
	"TukTuk/config"
	"regexp"
	"strconv"
	"strings"

	"github.com/aiomonitors/godiscord"
)

var SettingBot config.DiscordAlertSetting

func BotSendAlert(data, source_ip, time, ProtocolName string, id int64) {
	if SettingBot.Enabled {
		embed := godiscord.NewEmbed("TukTuk", "https://"+parseDomainforlink(config.Settings.DomainConfig.Name)+":1234/api/request/"+strings.ToLower(ProtocolName)+"?id="+strconv.Itoa(int(id)), config.Settings.DomainConfig.Name)
		embed.SetColor("F70505")

		embed.AddField("Received "+ProtocolName+" request from IP:"+source_ip, time, true)
		if SettingBot.LengthAlert == "Long" {
			embed.AddField("Data:", data, true)
		}
		err := embed.SendToWebhook(SettingBot.WebHook)
		if err != nil {
			panic(err)
		}
	}
}

func parseDomainforlink(str string) string {
	re := regexp.MustCompile("\\.")
	return re.Split(str, -1)[1] + "." + re.Split(str, -1)[2]
}
