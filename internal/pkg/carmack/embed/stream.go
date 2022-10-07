package embed

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed/color"
	"github.com/vikpe/qw-hub-api/types"
)

func FromStream(stream types.TwitchStream) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       stream.Channel,
		URL:         stream.Url,
		Description: stream.Title,
		Color:       color.Purple,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d viewers", stream.ViewerCount),
		},
	}

	return embed
}
