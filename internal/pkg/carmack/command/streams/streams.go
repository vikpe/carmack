package streams

import (
	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/carmack/internal/pkg/hub"
	"github.com/vikpe/carmack/internal/pkg/util"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "streams",
	Description: "Streams",
}

func Handler(i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	streams, err := hub.Streams()

	if 0 == len(streams) {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: util.StringOrError("No streams found.", err),
			},
		}
	}

	embeds := make([]*discordgo.MessageEmbed, 0)

	for _, stream := range streams {
		embeds = append(embeds, embed.FromStream(stream))
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  discordgo.MessageFlagsEphemeral,
			Embeds: embeds,
		},
	}
}
