package serverinfo

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/carmack/internal/pkg/discordbot"
	"github.com/vikpe/carmack/internal/pkg/util"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "server",
	Description: "Server info",
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
		genericServer, err := sstat.GetInfo(optionMap["address"].StringValue())

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}

		if err != nil {
			response.Data.Content = err.Error()
		}

		if genericServer.Version.IsMvdsv() {
			server := convert.ToMvdsv(genericServer)
			response.Data.Embeds = []*discordgo.MessageEmbed{embed.FromMvdsvServer(server)}
		} else if genericServer.Version.IsQtv() {
			server := convert.ToQtv(genericServer)
			response.Data.Embeds = []*discordgo.MessageEmbed{embed.FromQtvServer(server)}
		} else if genericServer.Version.IsQwfwd() {
			server := convert.ToQwfwd(genericServer)
			response.Data.Embeds = []*discordgo.MessageEmbed{embed.FromQwfwdServer(server)}
		} else {
			response.Data.Content = fmt.Sprintf("(server type not implemented: %s)", genericServer.Version.GetType())
		}

		return response
	}
}
