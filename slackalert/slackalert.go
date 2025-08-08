package slackalert

import (
	"TukTuk/config"

	"github.com/ashwanthkumar/slack-go-webhook"
)

func BotSendAlert(data string) {
	if config.Settings.SlackAlert.Enabled {
		paylod := slack.Payload{
			Text:     data,
			Username: config.Settings.SlackAlert.Username,
			Channel:  config.Settings.SlackAlert.Channel,
		}
		err := slack.Send(config.Settings.SlackAlert.WebHook, "", paylod)
		if len(err) > 0 {
			panic(err)
		}
	}
}
