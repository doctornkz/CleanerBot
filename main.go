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
	channelid string
	channelID int64
	response  apiResponse
)

const (
	translateURL = "https://translate.yandex.net/api/v1.5/tr.json/translate"
)

func init() {
	//yndAPIKeyPtr := flag.String("yndapikey", "", "YandexApiKey.")
	botAPIKeyPtr := flag.String("botapikey", "", "Bot ApiKey. See @BotFather messages for details")
	chatIDPtr := flag.String("chatid", "", "ChatID")
	channelIDPtr := flag.String("channelid", "", "ChannelID")

	flag.Parse()

	//yndAPIKey = *yndAPIKeyPtr
	botAPIKey = *botAPIKeyPtr
	chatid = *chatIDPtr
	channelid = *channelIDPtr
	if botAPIKey == "" || chatid == "" || channelid == "" {
		fmt.Println("Empty parameters, please use --help")
		os.Exit(1)
	}
	chatidInt32, err := strconv.Atoi(chatid)
	if err != nil {
		fmt.Println("Can't convert parameter to chatID, check input")
	}
	chatID = int64(chatidInt32)

	channelidInt32, err := strconv.Atoi(channelid)
	if err != nil {
		fmt.Println("Can't convert parameter to channelID, check input")
	}
	channelID = int64(channelidInt32)
}

func main() {

	bot, err := tgbotapi.NewBotAPI(botAPIKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("ID %d", bot.Self.ID)
	log.Printf("Is Bot? %v", bot.Self.IsBot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Chat == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("Message: %v", update.Message)
		log.Printf("ChatId: %v", update.Message.Chat.ID)

		isForward := false
		isMedia := false

		log.Println("Before check forward")

		if update.Message.ForwardFromChat != nil {
			if update.Message.ForwardFromChat.ID == channelID {
				log.Printf("forward %v", update.Message.ForwardFromChat.ID)
				isForward = true
			}
		}
		log.Println("After check forward")
		if update.Message.Photo != nil {
			isMedia = true
		}

		if update.Message.Video != nil {
			isMedia = true
		}

		log.Println("After check photo")
		if isForward && isMedia {
			log.Println("Forwarded photo, removing")
			deleteMessageConfig := tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: update.Message.MessageID}
			bot.DeleteMessage(deleteMessageConfig)
		} else {
			log.Println("Skipp message")
		}

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
