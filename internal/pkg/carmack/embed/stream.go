package embed

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed/color"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
)

func FromStream(stream twitch.Stream) *discordgo.MessageEmbed {
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
