package embed_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/assert/v2"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/carmack/internal/pkg/hub"
)

func TestFromStream(t *testing.T) {
	stream := hub.Stream{
		Channel:       "QuakeWorld",
		Url:           "https://twitch.tv/Quakeworld",
		Title:         "1on1: dough vs grl [ztndm3]",
		ViewerCount:   5,
		ServerAddress: "qw.foppa.dk:27502",
	}

	expect := discordgo.MessageEmbed{
		Title:       stream.Channel,
		URL:         stream.Url,
		Description: stream.Title,
		Color:       0xa970ff,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://static-cdn.jtvnw.net/previews-ttv/live_user_quakeworld-428x240.jpg",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Server",
				Value:  "`qw.foppa.dk:27502`",
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "5 viewers",
		},
	}

	assert.Equal(t, expect, embed.FromStream(stream))
}
