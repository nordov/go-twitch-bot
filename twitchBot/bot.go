package bot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/textproto"
	"regexp"
	"time"
)

var MsgRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)

type BotBody struct {
	Server  string
	Port    string
	Nick    string
	Token   string
	Channel string
	Conn    net.Conn
	MsgRate time.Duration
}

type Bot interface {
	Connect()
	Disconnect()
	JoinChannel()
	ListeningLoop()
	PostMsg()
}

type ChuckAPI struct {
	Value string `json:"value"`
}

func getChuckNorrisFactoid() string {

	resp, err := http.Get("https://api.chucknorris.io/jokes/random")

	if err != nil {
		fmt.Println("¯\\_(ツ)_/¯ Error found: Chuck Norris can't be bothered right now")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	factoid := ChuckAPI{}

	jsonErr := json.Unmarshal(body, &factoid)
	if jsonErr != nil {
		fmt.Println(err)
	}

	return factoid.Value

}

func log(str string) {

	timestamp := time.Now()

	fmt.Printf("[%s]: %s\n", timestamp.Format("01-02-2006 15:04:05"), str)

}

func (bot *BotBody) Connect() {

	var err error

	log("Attempting to connect to " + bot.Server)

	bot.Conn, err = net.Dial("tcp", bot.Server+":"+bot.Port)

	if err != nil {
		log("Failed connecting to " + bot.Server)

	} else {

		log("Connection with " + bot.Server + " stablished")

	}

}

func (bot *BotBody) JoinChannel() {

	log("Trying to join #" + bot.Channel)

	bot.Conn.Write([]byte("PASS " + bot.Token + "\r\n"))
	bot.Conn.Write([]byte("NICK " + bot.Nick + "\r\n"))
	bot.Conn.Write([]byte("JOIN #" + bot.Channel + "\r\n"))

	log("Joined #" + bot.Channel)

}

func (bot *BotBody) ListeningLoop() {

	tp := textproto.NewReader(bufio.NewReader(bot.Conn))

	for {
		line, err := tp.ReadLine()
		if err != nil {
			log("An error occurred fetching messages")
			break
		}

		if line == "PING :tmi.twitch.tv" {

			// respond to PING message with a PONG message, to maintain the connection
			bot.Conn.Write([]byte("PONG :tmi.twitch.tv\r\n"))
			continue
		}

		msgArray := MsgRegex.FindStringSubmatch(line)
		if msgArray != nil {
			msgUser := msgArray[1]
			msgType := msgArray[2]
			msgBody := msgArray[3]

			if msgType == "PRIVMSG" {
				log(msgUser + ": " + msgBody)
			}

			if msgBody == "!shutdown" {

				log("Shutdown requested by " + msgUser)
				log("Shutting down now...")
				bot.Disconnect()
				break

			}

			if msgBody == "!chucknorrisfact" {

				factoid := getChuckNorrisFactoid()
				log(factoid)

				bot.PostMsg(factoid)

			}

		}

		time.Sleep(bot.MsgRate)

	}

}

func (bot *BotBody) PostMsg(msg string) error {

	// check if message is too large
	if len(msg) > 512 {
		log("TWITCHBOT CONSOLE: Message is too big")
	}

	_, err := bot.Conn.Write([]byte("PRIVMSG #" + bot.Channel + " :" + msg + "\r\n"))
	if nil != err {
		log("TWITCHBOT CONSOLE: Error sending message")
	}
	return nil
}

func (bot *BotBody) Disconnect() {

	log("Disconnecting from " + bot.Server)

	bot.Conn.Close()

	log("Connection terminated")
}
