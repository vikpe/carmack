package embed

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/hub"
)

const colorPurple = 0xa970ff

func FromStream(stream hub.Stream) *discordgo.MessageEmbed {
	thumbSize := "428x240"
	thumbUrl := fmt.Sprintf(
		"https://static-cdn.jtvnw.net/previews-ttv/live_user_%s-%s.jpg",
		strings.ToLower(stream.Channel), thumbSize,
	)

	embed := &discordgo.MessageEmbed{
		Title:       stream.Channel,
		URL:         stream.Url,
		Description: stream.Title,
		Color:       colorPurple,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: thumbUrl},
		Fields:      []*discordgo.MessageEmbedField{},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d viewers", stream.ViewerCount),
		},
	}

	if stream.ServerAddress != "" {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:   "Server",
				Value:  fmt.Sprintf("`%s`", stream.ServerAddress),
				Inline: true,
			},
		)
	}

	return embed
}
