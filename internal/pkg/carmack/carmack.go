package carmack

import (
	"github.com/vikpe/carmack/internal/pkg/carmack/autocomplete"
	"github.com/vikpe/carmack/internal/pkg/carmack/command/demos"
	"github.com/vikpe/carmack/internal/pkg/carmack/command/findplayer"
	"github.com/vikpe/carmack/internal/pkg/carmack/command/serverinfo"
	"github.com/vikpe/carmack/internal/pkg/carmack/command/streams"
	"github.com/vikpe/carmack/internal/pkg/discordbot"
	"github.com/vikpe/serverstat"
)

type Carmack struct {
	discordbot.Bot
}

func New(token string, guildID string) (*Carmack, error) {
	bot, err := discordbot.New(token, guildID)

	statClient := serverstat.NewClient()
	bot.AddCommand(demos.Command, demos.Handler)
	bot.AddCommand(findplayer.Command, findplayer.Handler)
	bot.AddCommand(serverinfo.Command, serverinfo.GetHandler(statClient))
	bot.AddCommand(streams.Command, streams.GetHandler(statClient))
	bot.AddAutocompleteHandler("address", autocomplete.ServerAddress)

	return &Carmack{
		Bot: *bot,
	}, err
}
