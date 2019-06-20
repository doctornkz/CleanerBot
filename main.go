package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type apiResponse struct {
	StatusCode int      `json:"code"`
	Lang       string   `json:"lang"`
	Text       []string `json:"text"`
}

var (
	yndAPIKey string
	botAPIKey string
	chatid    string
	chatID    int64
	response  apiResponse
)

const (
	translateURL = "https://translate.yandex.net/api/v1.5/tr.json/translate"
)

func init() {
	//yndAPIKeyPtr := flag.String("yndapikey", "", "YandexApiKey.")
	botAPIKeyPtr := flag.String("botapikey", "", "Bot ApiKey. See @BotFather messages for details")
	ChatIDPtr := flag.String("chatid", "", "ChatID or ChannelID")

	flag.Parse()

	//yndAPIKey = *yndAPIKeyPtr
	botAPIKey = *botAPIKeyPtr
	chatid = *ChatIDPtr
	if botAPIKey == "" || chatid == "" {
		fmt.Println("Empty parameters, please use --help")
		os.Exit(1)
	}
	chatidInt32, err := strconv.Atoi(chatid)
	if err != nil {
		fmt.Println("Can't convert parameter to ChatID, check input")
	}
	chatID = int64(chatidInt32)

}

func main() {

	bot, err := tgbotapi.NewBotAPI(botAPIKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Chat == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("Message: %v", update.Message)
		log.Printf("ChatId: %v", update.Message.Chat.ID)
		// if photo if from channel
		//if update.Message.  {
		//	msg := tgbotapi.NewMessage(chatID, "Photo!")
		//}
		//bot.Send(msg)
	}
}

// Disabled
func getTranslate(phrase string) string {
	resp, err := http.PostForm(translateURL,
		url.Values{
			"key":  {yndAPIKey},
			"text": {phrase},
			"lang": {"ru-en"},
		})
	if err != nil {
		log.Fatal()
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal()
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Can't unmarshal")
		log.Fatal()
	}

	return response.Text[0]
}
