package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/discordbot"
	"github.com/vikpe/carmack/internal/pkg/hub"
	"github.com/vikpe/carmack/internal/pkg/util"
)

func ServerInfo() (*discordgo.ApplicationCommand, discordbot.CommandHandler) {
	cmd := discordgo.ApplicationCommand{
		Name:        "server",
		Description: "server command",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "address",
				Description:  "Server address",
				Required:     true,
				Autocomplete: true,
			},
		},
	}

	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		optionMap := util.ToOptionsMap(i.ApplicationCommandData().Options)
		server, err := hub.GetServerInfo(optionMap["address"].StringValue())
		responseContent := util.StringOrError(
			fmt.Sprintf("%s - %s", server.Address, server.Title),
			err,
		)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: responseContent,
			},
		})
	}
	return &cmd, handler
}
