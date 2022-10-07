package embed

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed/color"
	"github.com/vikpe/carmack/internal/pkg/util"
	"github.com/vikpe/qw-hub-api/types"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

func FromMvdsvServer(server mvdsv.Mvdsv) *discordgo.MessageEmbed {
	hostname := server.Settings.Get("hostname_parsed", server.Address)
	title := fmt.Sprintf(":flag_%s: %s", strings.ToLower(server.Geo.CC), hostname)
	statusText := strings.ToLower(fmt.Sprintf("%s - %s", server.Status.Name, server.Status.Description))

	clientFieldsInline := server.PlayerSlots.Used+server.SpectatorSlots.Used <= 6

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   fmt.Sprintf("Players (%d/%d)", server.PlayerSlots.Used, server.PlayerSlots.Total),
			Value:  util.SliceToNaturalList(analyze.GetPlayerPlainNames(server)),
			Inline: clientFieldsInline,
		},
		{
			Name:   fmt.Sprintf("Spectators (%d/%d)", server.SpectatorSlots.Used, server.SpectatorSlots.Total),
			Value:  util.SliceToNaturalList(server.SpectatorNames),
			Inline: clientFieldsInline,
		},
	}

	if len(server.QtvStream.Address) > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "QTV",
			Value: fmt.Sprintf("`%s`", server.QtvStream.Url),
		})
	}

	versionString := server.Settings.Get("*version", "")

	if server.Settings.Has("ktxver") {
		versionString += fmt.Sprintf(", KTX %s", server.Settings.Get("ktxver", "unknown"))
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf("%s on %s (%s)", server.Mode, server.Settings.Get("map", ""), statusText),
		Color:       color.Blue,
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Admin: %s [%s]", server.Settings.Get("*admin", "unknown"), versionString),
		},
	}

	return embed
}

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

func FromStream(stream types.TwitchStream) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       stream.Channel,
		URL:         stream.Url,
		Description: stream.Title,
		Color:       color.Purple,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d viewers", stream.ViewerCount),
		},
	}

	return embed
}
