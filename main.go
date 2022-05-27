package main

import (
	bot "go-twitch-bot/twitchBot"
	"time"
)

func main() {

	twitchBot := bot.BotBody{
		Server:  "irc.chat.twitch.tv",
		Port:    "6667",
		Nick:    "justinfan22",
		Token:   "oauth:yhr8n7e5hiqzspaxmgfzlecmvxej7b",
		Channel: "esfandtv",
		Conn:    nil,
		MsgRate: time.Duration(20/30) * time.Millisecond,
	}

	twitchBot.Connect()

	twitchBot.JoinChannel()

	twitchBot.ListeningLoop()

}
