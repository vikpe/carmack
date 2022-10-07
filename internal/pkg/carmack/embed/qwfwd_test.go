package embed_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

func TestFromQwfwdServer(t *testing.T) {
	server := qwfwd.Qwfwd{
		Address:     "46.227.68.148:30000",
		ClientNames: []string{"bps", "XantoM"},
		Settings: map[string]string{
			"*version":        "qwfwd 1.2",
			"hostname":        "QUAKE.SE KTX QWfwd",
			"hostname_parsed": "quake.se:30000",
			"maxclients":      "128",
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
		Title: ":flag_se: quake.se:30000",
		Color: 0x0c2aac,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Clients (2/128)",
				Value: "bps, XantoM",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "qwfwd 1.2",
		},
	}

	assert.Equal(t, expect, embed.FromQwfwdServer(server))
}
