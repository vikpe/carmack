package embed_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/qw-hub-api/types"
)

func TestFromStream(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		stream := types.TwitchStream{
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
	})

	t.Run("no server address", func(t *testing.T) {
		stream := types.TwitchStream{ServerAddress: ""}
		assert.Empty(t, embed.FromStream(stream).Fields)
	})
}
