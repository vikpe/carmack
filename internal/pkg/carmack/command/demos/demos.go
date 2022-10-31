package demos

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/util"
	"github.com/vikpe/go-qwhub"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "demos",
	Description: "Find demos",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "query",
			Description: "space separated list of player/team names and map. example: lordlame dm4",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "mode",
			Description: "name of mode (duel, 2on2, 4on4)",
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "qtv_server",
			Description: "qtv server (qw.foppa.dk)",
		},
	},
}

func Handler(i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	optionMap := util.ToOptionsMap(i.ApplicationCommandData().Options)

	params := map[string]string{
		"q":          optionMap["query"].StringValue(),
		"mode":       optionMap["mode"].StringValue(),
		"qtv_server": optionMap["qtv_server"].StringValue(),
		"limit":      "16",
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
	header := fmt.Sprintf(`**Demo search**: found %d demo(s)`, len(demos))
	demoCount := len(demos)

	if 0 == demoCount {
		return header
	}

	if demoCount > 15 {
		header += " (showing 15 most recent)"
	}

	responseLines := []string{header}

	for _, demo := range demos {
		/*timestamp := demo.Time.Format("2006-01-02 15:04")
		downloadLink := fmt.Sprintf("[%s](%s)", strings.ReplaceAll(demo.Filename, "_", " "), demo.DownloadUrl)
		line := fmt.Sprintf("`%s` - %s", timestamp, downloadLink)*/
		responseLines = append(responseLines, demo.Filename)
	}

	return strings.Join(responseLines, "\n")
}
