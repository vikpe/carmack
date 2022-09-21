package embed

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/hub"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
)

const colorPurple = 0xa970ff
const colorBlue = 0x0c2aac

func FromMvdsvServer(server mvdsv.Mvdsv) *discordgo.MessageEmbed {
	hostname := server.Settings.Get("hostname_parsed", server.Address)
	title := fmt.Sprintf(":flag_%s: %s", strings.ToLower(server.Geo.CC), hostname)

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf("%s on %s\n", server.Mode, server.Settings.Get("map", "")),
		Color:       colorBlue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("Players (%d/%d)", server.PlayerSlots.Used, server.PlayerSlots.Total),
				Value: strings.Join(analyze.GetPlayerPlainNames(server), ", "),
			},
			{
				Name:  fmt.Sprintf("Spectators (%d/%d)", server.SpectatorSlots.Used, server.SpectatorSlots.Total),
				Value: strings.Join(server.SpectatorNames, ", "),
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

func FromStream(stream hub.Stream) *discordgo.MessageEmbed {
	thumbSize := "428x240"
	thumbUrl := fmt.Sprintf(
		"https://static-cdn.jtvnw.net/previews-ttv/live_user_%s-%s.jpg",
		strings.ToLower(stream.Channel), thumbSize,
	)

	embed := &discordgo.MessageEmbed{
		Title:       stream.Channel,
		URL:         stream.Url,
		Description: stream.Title,
		Color:       colorPurple,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: thumbUrl},
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
