package embed

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed/color"
	"github.com/vikpe/carmack/internal/pkg/util"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/qclient"
)

func FromGenericServer(server qserver.GenericServer) *discordgo.MessageEmbed {
	hostname := server.Settings.Get("hostname_parsed", server.Address)
	title := fmt.Sprintf(":flag_%s: %s", strings.ToLower(server.Geo.CC), hostname)

	admin := "unknown"

	if server.Settings.Has("admin") {
		admin = server.Settings.Get("admin", admin)
	} else if server.Settings.Has("*admin") {
		admin = server.Settings.Get("*admin", admin)
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: server.Settings.Get("hostname", "-"),
		Color:       color.Blue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Players",
				Value: util.SliceToNaturalList(qclient.ClientNames(server.Players())),
			},
			{
				Name:  "Spectators",
				Value: util.SliceToNaturalList(qclient.ClientNames(server.Spectators())),
			},
			{
				Name:   "Gamedir",
				Value:  server.Settings.Get("*gamedir", "(unknown)"),
				Inline: true,
			},
			{
				Name:   "Map",
				Value:  server.Settings.Get("map", "(unknown)"),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf(
				"Admin: %s [%s]",
				admin,
				server.Settings.Get("*version", "unknown version"),
			),
		},
	}

	return embed
}
