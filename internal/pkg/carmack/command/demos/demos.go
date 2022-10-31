package demos

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/util"
	"github.com/vikpe/go-qwhub"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "demos",
	Description: "Find demos",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "q",
			Description: `query: space separated list of players/teams/map. example: "xantom dm4"`,
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "mode",
			Description: "Game mode",
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: "duel", Value: "duel"},
				{Name: "2on2", Value: "2on2"},
				{Name: "4on4", Value: "4on4"},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "qtv_server",
			Description: "QTV server",
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: "dal.spawnfrag.com", Value: "dal.spawnfrag.com:28000"},
				{Name: "de.aye.wtf", Value: "de.aye.wtf:28000"},
				{Name: "la.spawnfrag.com", Value: "la.spawnfrag.com:28000"},
				{Name: "london.badplace.eu", Value: "london.badplace.eu:28000"},
				{Name: "nl.aye.wtf", Value: "nl.aye.wtf:28000"},
				{Name: "ny.qwsrv.com", Value: "ny.qwsrv.com:28000"},
				{Name: "play.quake1.pl", Value: "play.quake1.pl:28000"},
				{Name: "qtv.nicotinelounge.com", Value: "qtv.nicotinelounge.com:443"},
				{Name: "quake.se", Value: "quake.se:28000"},
				{Name: "qw.foppa.dk", Value: "qw.foppa.dk:28000"},
				{Name: "qw.irc.ax", Value: "qw.irc.ax:28000"},
				{Name: "troopers.fi", Value: "troopers.fi:28000"},
			},
		},
	},
}

func Handler(i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	optionMap := util.ToOptionsMap(i.ApplicationCommandData().Options)
	params := map[string]string{"limit": "11"}

	for _, opt := range Command.Options {
		if _, ok := optionMap[opt.Name]; ok {
			params[opt.Name] = optionMap[opt.Name].StringValue()
		}
	}

	demos := qwhub.NewClient().Demos(params)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: ContentFromDemos(demos),
		},
	}

	return response
}

func ContentFromDemos(demos []qtvscraper.Demo) string {
	header := "**Demo search**: "
	demoCount := len(demos)

	if 0 == demoCount {
		header += " found 0 demos."
		return header
	}

	if demoCount > 10 {
		header += " found 10+ demos (showing 10 most recent)"
		demos = demos[0:10]
	} else {
		header += fmt.Sprintf("found %d demo(s)", demoCount)
	}

	responseLines := []string{header}

	for _, demo := range demos {
		timestamp := demo.Time.Format("2006-01-02 15:04")
		demoFilename := qdemo.Filename(demo.Filename)
		participants := formatParticipants(demoFilename.Participants())
		linkText := fmt.Sprintf("%s - %s \\[%s\\]", demoFilename.Mode(), participants, demoFilename.Map())
		line := fmt.Sprintf("`%s` - [%s](%s)", timestamp, linkText, demo.DownloadUrl)
		responseLines = append(responseLines, line)
	}

	return strings.Join(responseLines, "\n")
}

func formatParticipants(participants []string) string {
	parts := make([]string, 0)

	for _, p := range participants {
		parts = append(parts, strings.Trim(p, "][()"))
	}

	return strings.Join(parts, " vs ")
}
