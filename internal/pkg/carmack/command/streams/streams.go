package streams

import (
	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/go-qwhub"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "streams",
	Description: "Streams",
}

func Handler(i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
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
		embeds = append(embeds, embed.FromStream(stream))
	}

	response.Data.Embeds = embeds
	return response
}
