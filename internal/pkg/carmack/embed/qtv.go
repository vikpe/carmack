package embed

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed/color"
	"github.com/vikpe/carmack/internal/pkg/util"
	"github.com/vikpe/serverstat/qserver/qtv"
)

func FromQtvServer(server qtv.Qtv) *discordgo.MessageEmbed {
	hostname := server.Settings.Get("hostname_parsed", server.Address)
	title := fmt.Sprintf(":flag_%s: %s", strings.ToLower(server.Geo.CC), hostname)

	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: color.Blue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("Spectators (%d/%d)", len(server.SpectatorNames), server.Settings.GetInt("maxclients", 0)),
				Value: util.SliceToNaturalList(server.SpectatorNames),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: server.Settings.Get("*version", "(unknown version)"),
		},
	}

	return embed
}
