package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vikpe/carmack/internal/pkg/carmack"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	carmackBot, err := carmack.New(
		os.Getenv("BOT_TOKEN"),
		os.Getenv("GUILD_ID"),
	)

	if err != nil {
		log.Fatal("unable to create carmackBot", err)
		return
	}

	carmackBot.Start() // blocking operation
}
