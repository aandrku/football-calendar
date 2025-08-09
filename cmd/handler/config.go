package main

import "os"

type Config struct {
	endpoint         string
	apiKey           string
	botToken         string
	chatID           string
	telegramEndpoint string
}

func LoadConfig() Config {
	var c Config

	c.endpoint = os.Getenv("API_ENDPOINT")
	c.apiKey = os.Getenv("API_KEY")
	c.botToken = os.Getenv("BOT_TOKEN")
	c.chatID = os.Getenv("CHAT_ID")
	c.telegramEndpoint = "https://api.telegram.org/bot%s/sendMessage"

	return c
}
