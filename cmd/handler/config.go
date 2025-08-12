package main

import (
	"net/http"
	"os"
)

type Config struct {
	endpoint         string
	apiKey           string
	botToken         string
	chatID           string
	telegramEndpoint string
	leaguesFile      string
	client           *http.Client
	season           string
	timezone         string
}

func LoadConfig() Config {
	var c Config

	c.endpoint = os.Getenv("API_ENDPOINT")
	c.apiKey = os.Getenv("API_KEY")
	c.botToken = os.Getenv("BOT_TOKEN")
	c.chatID = os.Getenv("CHAT_ID")
	c.telegramEndpoint = "https://api.telegram.org/bot%s/sendMessage"
	c.leaguesFile = "./leagues.json"
	c.client = &http.Client{}
	c.season = os.Getenv("SEASON")
	c.timezone = os.Getenv("TIMEZONE")

	return c
}
