package embed

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed/color"
	"github.com/vikpe/carmack/internal/pkg/util"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

func FromQwfwdServer(server qwfwd.Qwfwd) *discordgo.MessageEmbed {
	hostname := server.Settings.Get("hostname_parsed", server.Address)
	title := fmt.Sprintf(":flag_%s: %s", strings.ToLower(server.Geo.CC), hostname)

	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: color.Blue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("Clients (%d/%d)", len(server.ClientNames), server.Settings.GetInt("maxclients", 0)),
				Value: util.SliceToNaturalList(server.ClientNames),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: server.Settings.Get("*version", "(unknown version)"),
		},
	}

	return embed
}
