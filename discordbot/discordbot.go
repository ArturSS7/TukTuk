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
		//Print Name author and link to the admin panel
		embed := godiscord.NewEmbed("TukTuk", "https://"+parseDomainforlink(config.Settings.DomainConfig.Name)+":1234/api/request/"+strings.ToLower(ProtocolName)+"?id="+strconv.Itoa(int(id)), config.Settings.DomainConfig.Name)
		//Color - red(maybe add to the config file )
		embed.SetColor("F70505")
		//Short format - Protocol name, source_ip and time
		embed.AddField("Received "+ProtocolName+" request from IP:"+source_ip, time, true)
		//Long format - same + data
		if SettingBot.LengthAlert == "Long" {
			embed.AddField("Data:", data, true)
		}
		//Send
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
