package embed_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestFromGenericServer(t *testing.T) {
	server := qserver.GenericServer{
		Address: "177.71.158.196:27500",
		Settings: map[string]string{
			"*gamedir":        "fortress",
			"*version":        "FTE SVN 5468",
			"admin":           "Mouse",
			"hostname":        "#00 HueTF",
			"hostname_parsed": "177.71.158.196:27500",
			"map":             "wellgl1",
		},
		Clients: []qclient.Client{{
			Name:   qstring.New("XantoM"),
			Team:   qstring.New("red"),
			Skin:   "",
			Colors: [2]uint8{13, 13},
			Frags:  2,
			Ping:   38,
			Time:   4,
		}},
		Geo: geo.Location{
			CC:          "BR",
			Country:     "Brazil",
			Region:      "South America",
			City:        "São Paulo",
			Coordinates: [2]float32{-23.5505, -46.6333},
		},
	}

	expect := &discordgo.MessageEmbed{
		Title:       ":flag_br: 177.71.158.196:27500",
		Description: "#00 HueTF",
		Color:       0x0c2aac,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Players",
				Value: "XantoM",
			},
			{
				Name:  "Spectators",
				Value: "-",
			},
			{
				Name:   "Gamedir",
				Value:  "fortress",
				Inline: true,
			},
			{
				Name:   "Map",
				Value:  "wellgl1",
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Admin: Mouse [FTE SVN 5468]",
		},
	}

	assert.Equal(t, expect, embed.FromGenericServer(server))
}
