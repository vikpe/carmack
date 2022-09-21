package serverinfo

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/discordbot"
	"github.com/vikpe/carmack/internal/pkg/util"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
)

var Command = &discordgo.ApplicationCommand{
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

func GetHandler(sstat *serverstat.Client) discordbot.CommandHandler {
	return func(i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
		optionMap := util.ToOptionsMap(i.ApplicationCommandData().Options)
		serverInfo, err := sstat.GetInfo(optionMap["address"].StringValue())
		responseContent := ""

		if serverInfo.Version.IsMvdsv() {
			server := convert.ToMvdsv(serverInfo)
			responseContent = fmt.Sprintf("%s - %s", server.Address, server.Title)
		}

		responseContent = util.StringOrError(responseContent, err)

		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: responseContent,
			},
		}
	}
}
