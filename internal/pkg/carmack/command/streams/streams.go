package streams

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
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

	components := make([]discordgo.MessageComponent, 0)

	for _, s := range streams {
		label := fmt.Sprintf("%s - %s (%d viewers)\n", s.Channel, s.Title, s.ViewerCount)
		components = append(components,
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Style: discordgo.LinkButton,
						Label: label,
						URL:   s.Url,
					},
				},
			})
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: components,
		},
	}
}
