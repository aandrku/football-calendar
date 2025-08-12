package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"

	"errors"
	"fmt"
	"os"
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

	leagues, err := readLeagues(c)
	if err != nil {
		log.Fatalf("There was an unrecoverable error: %v", err)
	}

	m, err := getModels(leagues, c)
	if err != nil {
		log.Fatalf("There was an unrecoverable error: %v", err)

	}

	err = notifyUser(m, c)
	if err != nil {
		log.Fatalf("There was an error while sending a text to user %v", err)
	}

}

type League struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func readLeagues(config Config) ([]League, error) {
	var leagues []League

	f, err := os.Open(config.leaguesFile)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)

	err = dec.Decode(&leagues)
	if err != nil {
		return nil, err
	}

	return leagues, nil

}

func getModels(leagues []League, config Config) ([]Model, error) {
	c := config.client

	models := make([]Model, 0, len(leagues))
	var lock sync.Mutex
	var wg sync.WaitGroup

	for _, l := range leagues {
		wg.Add(1)
		go func() {
			req, err := http.NewRequest("GET", config.endpoint, nil)
			if err != nil {
				log.Printf("error creating a request while fetching %q league data %v", l.Name, err)
			}

			req.Header.Add("x-rapidapi-key", config.apiKey)

			q := req.URL.Query()

			q.Set("league", l.ID)
			q.Set("date", time.Now().Format("2006-01-02"))
			q.Set("season", config.season)
			q.Set("timezone", config.timezone)

			// forgot to set the actual query
			req.URL.RawQuery = q.Encode()

			res, err := c.Do(req)
			if err != nil {
				log.Printf("error making a request while fetching %q league data %v", l.Name, err)
			}

			dec := json.NewDecoder(res.Body)

			var r ApiResponse
			err = dec.Decode(&r)
			if err != nil {
				log.Printf("error decoding json while fetching %q league data %v", l.Name, err)
			}

			m := Model{
				League: l.Name,
				Games:  make([]GameModel, 0),
			}

			for _, v := range r.Response {
				var game GameModel

				game.Home = v.Teams.Home.Name
				game.Away = v.Teams.Away.Name
				game.Time = v.Fixture.Date

				m.Games = append(m.Games, game)
			}

			lock.Lock()
			models = append(models, m)
			lock.Unlock()

			wg.Done()
		}()

	}

	wg.Wait()

	return models, nil
}

func notifyUser(models []Model, config Config) error {
	type RequestBody struct {
		ChatID    string `json:"chat_id"`
		Text      string `json:"text"`
		ParseMode string `json:"parse_mode"`
	}

	// build message
	b := strings.Builder{}

	for _, m := range models {
		b.WriteString(m.Format())
		b.WriteRune('\n')
	}

	url := fmt.Sprintf(config.telegramEndpoint, config.botToken)

	rb := RequestBody{
		ChatID:    config.chatID,
		Text:      b.String(),
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
