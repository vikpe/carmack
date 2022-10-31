package embed_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
)

func TestFromStream(t *testing.T) {
	stream := twitch.Stream{
		Channel:       "QuakeWorld",
		Url:           "https://twitch.tv/Quakeworld",
		Title:         "1on1: dough vs grl [ztndm3]",
		ViewerCount:   5,
		ServerAddress: "qw.foppa.dk:27502",
	}

	expect := &discordgo.MessageEmbed{
		Title:       "QuakeWorld",
		URL:         "https://twitch.tv/Quakeworld",
		Description: "1on1: dough vs grl [ztndm3]",
		Color:       0xa970ff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "5 viewers",
		},
	}

	assert.Equal(t, expect, embed.FromStream(stream))
}
