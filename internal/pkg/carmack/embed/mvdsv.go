package embed

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed/color"
	"github.com/vikpe/carmack/internal/pkg/util"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
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
