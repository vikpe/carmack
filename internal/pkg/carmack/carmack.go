package carmack

import (
	"github.com/vikpe/carmack/internal/pkg/carmack/autocomplete"
	"github.com/vikpe/carmack/internal/pkg/carmack/command"
	"github.com/vikpe/carmack/internal/pkg/discordbot"
)

type Carmack struct {
	discordbot.Bot
}

func New(token string, guildID string) (*Carmack, error) {
	bot, err := discordbot.New(token, guildID)

	bot.AddCommand(command.ServerInfo())
	bot.AddCommand(command.FindPlayer())
	bot.AddAutocompleteHandler("address", autocomplete.ServerAddress)

	return &Carmack{Bot: *bot}, err
}
