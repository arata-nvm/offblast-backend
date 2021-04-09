package domain

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/arata-nvm/offblast/config"
)

func PostToSlack(body interface{}) {
	text, err := json.Marshal(body)
	if err != nil {
		log.Fatalln(err)
	}

	payload, err := json.Marshal(map[string]interface{}{"text": string(text)})
	if err != nil {
		log.Fatalln(err)
	}

	url := config.SlackWebhook()
	http.Post(url, "application/json", bytes.NewBuffer(payload))
}
