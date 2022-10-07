package streams

import (
	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed/color"
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
		streamCount := len(streams)

		if 0 == streamCount {
			response.Data.Content = "No streams found."
			return response
		}

		embeds := make([]*discordgo.MessageEmbed, 0)

		for index, stream := range streams {
			embedColor := streamEmbedColor(index, streamCount)

			streamEmbed := embed.FromStream(stream)
			streamEmbed.Color = embedColor
			embeds = append(embeds, streamEmbed)

			if len(stream.ServerAddress) > 0 {
				genericServer, err := sstat.GetInfo(stream.ServerAddress)

				if err == nil && genericServer.Version.IsMvdsv() {
					mvdsvServer := convert.ToMvdsv(genericServer)
					serverEmbed := embed.FromMvdsvServer(mvdsvServer)
					serverEmbed.Color = embedColor
					embeds = append(embeds, serverEmbed)
				}
			}
		}

		response.Data.Embeds = embeds
		return response
	}
}

func streamEmbedColor(index int, total int) int {
	if 1 == total {
		return color.Purple
	}

	return color.FromIndex(index)
}
