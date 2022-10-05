package streams

import (
	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/carmack/internal/pkg/discordbot"
	"github.com/vikpe/go-qwhub"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "streams",
	Description: "Streams",
}

func GetHandler(sstat *serverstat.Client) discordbot.CommandHandler {
	return func(i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}

		streams := qwhub.NewClient().Streams()

		if 0 == len(streams) {
			response.Data.Content = "No streams found."
			return response
		}

		embeds := make([]*discordgo.MessageEmbed, 0)

		for _, stream := range streams {
			streamEmbed := embed.FromStream(stream)
			embeds = append(embeds, streamEmbed)

			if len(stream.ServerAddress) > 0 {
				genericServer, err := sstat.GetInfo(stream.ServerAddress)

				if err == nil && genericServer.Version.IsMvdsv() {
					mvdsvServer := convert.ToMvdsv(genericServer)
					serverEmbed := embed.FromMvdsvServer(mvdsvServer)
					embeds = append(embeds, serverEmbed)
				}
			}
		}

		response.Data.Embeds = embeds
		return response
	}
}
