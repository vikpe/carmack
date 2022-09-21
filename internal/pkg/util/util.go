package util

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ToOptionsMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

func StringOrError(str string, err error) string {
	if err != nil {
		return err.Error()
	}

	return str
}

func ContainsAllSubstrings(haystack string, needles []string) bool {
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			return false
		}
	}

	return true
}

func ToOptionChoices(options [][]string) []*discordgo.ApplicationCommandOptionChoice {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0)

	for _, opt := range options {
		choices = append(choices, ToOptionsChoice(opt))
	}

	return choices
}

func ToOptionsChoice(option []string) *discordgo.ApplicationCommandOptionChoice {
	choice := &discordgo.ApplicationCommandOptionChoice{Name: "", Value: ""}
	optionLen := len(option)

	if 0 == optionLen {
		return choice
	}

	choice.Name = option[0]

	if 1 == optionLen {
		choice.Value = choice.Name
	} else {
		choice.Value = option[1]
	}

	return choice
}
