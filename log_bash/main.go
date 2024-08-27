package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"net/http"
	"net/url"
)

const tgAPIToken = "861103068:AAGMGOzTCmIhnALCOQzs9H0fAPQknwoAX9s"
const tgUserID = "547603956"

func sendToTelegram(message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tgAPIToken)
	data := url.Values{}
	data.Set("chat_id", tgUserID)
	data.Set("text", message)

	_, err := http.PostForm(apiURL, data)
	if err != nil {
		log.Printf("Failed to send message to Telegram: %v\n", err)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		sendToTelegram(input)
	}
}
