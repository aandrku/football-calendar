package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"strings"
	"sync"

	"errors"
	"fmt"
	"time"

	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	c := LoadConfig()
	fn := function{
		config: c,
	}

	lambda.Start(fn.handleRequest)
}

type League struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// function provides methods that implement AWS lambda function functionality.
//
// Before being used, function must be loaded with configuration.
type function struct {
	config Config
}

func (fn function) handleRequest() {
	leagues, err := fn.readLeagues()
	if err != nil {
		log.Fatalf("There was an unrecoverable error: %v", err)
	}

	m, err := fn.getModels(leagues)
	if err != nil {
		log.Fatalf("There was an unrecoverable error: %v", err)

	}

	err = fn.notifyUser(m)
	if err != nil {
		log.Fatalf("There was an error while sending a text to user %v", err)
	}

}

//go:embed leagues.json
var leaguesJSON []byte

func (fn function) readLeagues() ([]League, error) {
	var leagues []League

	r := bytes.NewReader(leaguesJSON)

	dec := json.NewDecoder(r)

	err := dec.Decode(&leagues)
	if err != nil {
		return nil, err
	}

	return leagues, nil

}

func (fn function) getModels(leagues []League) ([]Model, error) {
	c := fn.config.client

	models := make([]Model, 0, len(leagues))
	var lock sync.Mutex
	var wg sync.WaitGroup

	for _, l := range leagues {
		wg.Add(1)
		go func() {
			req, err := http.NewRequest("GET", fn.config.endpoint, nil)
			if err != nil {
				log.Printf("error creating a request while fetching %q league data %v", l.Name, err)
			}

			req.Header.Add("x-rapidapi-key", fn.config.apiKey)

			q := req.URL.Query()

			q.Set("league", l.ID)
			q.Set("date", time.Now().Format("2006-01-02"))
			q.Set("season", fn.config.season)
			q.Set("timezone", fn.config.timezone)

			// forgot to set the actual query
			req.URL.RawQuery = q.Encode()

			res, err := c.Do(req)
			if err != nil {
				log.Printf("error making a request while fetching %q league data %v", l.Name, err)
			}
			defer res.Body.Close()

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

func (fn function) notifyUser(models []Model) error {
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

	url := fmt.Sprintf(fn.config.telegramEndpoint, fn.config.botToken)

	rb := RequestBody{
		ChatID:    fn.config.chatID,
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
