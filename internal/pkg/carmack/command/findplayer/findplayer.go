package findplayer

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/hub"
	"github.com/vikpe/carmack/internal/pkg/util"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "find",
	Description: "Find player",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "Player name",
			Required:    true,
		},
	},
}

func Handler(i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	optionMap := util.ToOptionsMap(i.ApplicationCommandData().Options)
	playerName := optionMap["name"].StringValue()
	server, err := hub.FindPlayerOnServer(playerName)
	responseContent := util.StringOrError(
		fmt.Sprintf("%s is playing at %s (%s)", playerName, server.Address, server.Title),
		err,
	)

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: responseContent,
		},
	}
}
