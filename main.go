// e*;6J.Erbv7/G^*
// oauth:yhr8n7e5hiqzspaxmgfzlecmvxej7b

// Client ID: pti5lzx6pn2z46ip9wrsjid1wcy5mg
package main

import (
	"time"

	"go-twitch-bot/bot"
)

func main() {

	twitchBot := bot.BotBody{
		server:  "irc.chat.twitch.tv",
		port:    "6667",
		nick:    "justinfan22",
		token:   "oauth:yhr8n7e5hiqzspaxmgfzlecmvxej7b",
		channel: "esfandtv",
		conn:    nil,
		MsgRate: time.Duration(20/30) * time.Millisecond,
	}

	twitchBot.Connect()

	// twitchBot.joinChannel()

}
