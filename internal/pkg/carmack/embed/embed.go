package embed

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/qw-hub-api/types"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
)

const colorPurple = 0xa970ff
const colorBlue = 0x0c2aac

func sliceToNaturalList(values []string) string {
	if 0 == len(values) {
		return "-"
	}

	return strings.Join(values, ", ")
}

func FromMvdsvServer(server mvdsv.Mvdsv) *discordgo.MessageEmbed {
	hostname := server.Settings.Get("hostname_parsed", server.Address)
	title := fmt.Sprintf(":flag_%s: %s", strings.ToLower(server.Geo.CC), hostname)

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf("%s on %s", server.Mode, server.Settings.Get("map", "")),
		Color:       colorBlue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("Players (%d/%d)", server.PlayerSlots.Used, server.PlayerSlots.Total),
				Value: sliceToNaturalList(analyze.GetPlayerPlainNames(server)),
			},
			{
				Name:  fmt.Sprintf("Spectators (%d/%d)", server.SpectatorSlots.Used, server.SpectatorSlots.Total),
				Value: sliceToNaturalList(server.SpectatorNames),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf(
				"%s - %s",
				server.Status.Name, server.Status.Description,
			),
		},
	}

	return embed
}

func FromStream(stream types.TwitchStream) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       stream.Channel,
		URL:         stream.Url,
		Description: stream.Title,
		Color:       colorPurple,
		Fields:      []*discordgo.MessageEmbedField{},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d viewers", stream.ViewerCount),
		},
	}

	if stream.ServerAddress != "" {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:   "Server",
				Value:  fmt.Sprintf("`%s`", stream.ServerAddress),
				Inline: true,
			},
		)
	}

	return embed
}
