package streams

import (
	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/carmack/internal/pkg/hub"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "streams",
	Description: "Streams",
}

func Handler(i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	streams, err := hub.NewClient().Streams()

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}

	if err != nil {
		response.Data.Content = err.Error()
		return response
	}

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
