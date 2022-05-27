package bot

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"regexp"
	"time"
)

var MsgRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)

type BotBody struct {
	server  string
	port    string
	nick    string
	token   string
	channel string
	conn    net.Conn
	MsgRate time.Duration
}

type Bot interface {
	Connect()
	Disconnect()
	joinChannel()
	listeningLoop()
}

func log(str string) {

	timestamp := time.Now()

	fmt.Printf("[%s]: %s\n", timestamp.Format("01-02-2006 15:04:05"), str)

}

func (bot *Bot) Connect() {

	var err error

	log("Attempting to connect to " + bot.server)

	bot.conn, err = net.Dial("tcp", bot.server+":"+bot.port)

	if err != nil {
		log("Failed connecting to " + bot.server)

	} else {

		log("Connection with " + bot.server + " stablished")

	}

}

func (bot *Bot) joinChannel() {

	log("Trying to join #" + bot.channel)

	bot.conn.Write([]byte("PASS " + bot.token + "\r\n"))
	bot.conn.Write([]byte("NICK " + bot.nick + "\r\n"))
	bot.conn.Write([]byte("JOIN #" + bot.channel + "\r\n"))

	log("Joined #" + bot.channel)

}

func (bot *Bot) listeningLoop() {

	tp := textproto.NewReader(bufio.NewReader(bot.conn))

	for {
		line, err := tp.ReadLine()
		if err != nil {
			log("An error occurred fetching messages")
			break
		}

		log(line)

		if "PING :tmi.twitch.tv" == line {

			// respond to PING message with a PONG message, to maintain the connection
			bot.conn.Write([]byte("PONG :tmi.twitch.tv\r\n"))
			continue
		}

		// handle a PRIVMSG message
		matches := MsgRegex.FindStringSubmatch(line)
		if matches != nil {

			userName := matches[1]
			msgType := matches[2]

			if msgType == "PRIVMSG" {
				msg := matches[3]
				log(userName + ": " + msg + "\n")

				switch msg {
				case "!shutdown":
					log("Shutdown command received. Shutting down now...")
					bot.Disconnect()
					break
				case "!chucknorris":
					log("CHUCK NORRIS FACTOID GOES HERE")
				default:
					// do nothing
				}

			}
		}

		time.Sleep(bot.MsgRate)

	}

}

func (bot *Bot) Disconnect() {

	log("Disconnecting from " + bot.server)

	bot.conn.Close()

	log("Connection terminated")
}
