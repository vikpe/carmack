package embed_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qtv"
)

func TestFromQtvServer(t *testing.T) {
	server := qtv.Qtv{
		Address:        "46.227.68.148:28000",
		SpectatorNames: []string{"bps", "XantoM"},
		Settings: map[string]string{
			"*version":        "QTV 1.12-rc1",
			"hostname":        "QUAKE.SE KTX Qtv",
			"hostname_parsed": "quake.se:28000",
			"maxclients":      "100",
		},
		Geo: geo.Location{
			CC:          "SE",
			Country:     "Sweden",
			Region:      "Europe",
			City:        "Hagersten",
			Coordinates: [2]float32{59.2885, 17.9612},
		},
	}

	expect := &discordgo.MessageEmbed{
		Title: ":flag_se: quake.se:28000",
		Color: 0x0c2aac,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Spectators (2/100)",
				Value: "bps, XantoM",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "QTV 1.12-rc1",
		},
	}

	assert.Equal(t, expect, embed.FromQtvServer(server))
}
