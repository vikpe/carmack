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
