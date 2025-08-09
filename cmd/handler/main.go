package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("There was an unrecoverable error: %v", err)
	}

	c := LoadConfig()

	m, err := getModel(c)
	if err != nil {
		log.Fatalf("There was an unrecoverable error: %v", err)

	}

	err = notifyUser(m, c)
	if err != nil {
		log.Fatalf("There was an error while sending a text to user %v", err)
	}

}

func getModel(config Config) (Model, error) {
	req, err := http.NewRequest("GET", config.endpoint, nil)
	if err != nil {
		return Model{}, err
	}

	req.Header.Add("x-rapidapi-key", config.apiKey)

	q := req.URL.Query()

	q.Set("league", "253")
	q.Set("date", time.Now().Format("2006-01-02"))
	q.Set("season", "2025")
	q.Set("timezone", "America/New_York")

	// forgot to set the actual query
	req.URL.RawQuery = q.Encode()

	// we need a client to make the request
	c := http.Client{}

	res, err := c.Do(req)
	if err != nil {
		return Model{}, err
	}

	dec := json.NewDecoder(res.Body)

	var r ApiResponse
	err = dec.Decode(&r)
	if err != nil {
		return Model{}, err
	}

	m := Model{
		League: "Europa League",
		Games:  make([]GameModel, 0),
	}

	for _, v := range r.Response {
		var game GameModel

		game.Home = v.Teams.Home.Name
		game.Away = v.Teams.Away.Name
		game.Time = v.Fixture.Date

		m.Games = append(m.Games, game)
	}

	return m, nil
}

func notifyUser(model Model, config Config) error {
	type RequestBody struct {
		ChatID    string `json:"chat_id"`
		Text      string `json:"text"`
		ParseMode string `json:"parse_mode"`
	}

	url := fmt.Sprintf(config.telegramEndpoint, config.botToken)

	rb := RequestBody{
		ChatID:    config.chatID,
		Text:      model.Format(),
		ParseMode: "Markdown",
	}

	data, err := json.Marshal(&rb)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("message not sent")
	}

	return nil
}
