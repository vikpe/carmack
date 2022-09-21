package autocomplete_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/assert/v2"
	"github.com/vikpe/carmack/internal/pkg/carmack/autocomplete"
)

func getDataOption(value string) *discordgo.ApplicationCommandInteractionDataOption {
	return &discordgo.ApplicationCommandInteractionDataOption{
		Name:    "address",
		Type:    discordgo.ApplicationCommandOptionString,
		Value:   value,
		Focused: true,
	}
}

func TestServerAddress(t *testing.T) {
	t.Run("single needle", func(t *testing.T) {
		dataOption := getDataOption("foppa")
		choices := autocomplete.ServerAddress(dataOption)
		expect := [][]string{
			{"foppa", "choice_custom"},
			{"qw.foppa.dk:27500", "qw.foppa.dk:27500"},
			{"qw.foppa.dk:27501", "qw.foppa.dk:27501"},
			{"qw.foppa.dk:27502", "qw.foppa.dk:27502"},
			{"qw.foppa.dk:27503", "qw.foppa.dk:27503"},
			{"qw.foppa.dk:27504", "qw.foppa.dk:27504"},
			{"qw.foppa.dk:27505", "qw.foppa.dk:27505"},
		}
		assert.Equal(t, expect, choices)
	})

	t.Run("multiple needles", func(t *testing.T) {
		dataOption := getDataOption("foppa 502")
		choices := autocomplete.ServerAddress(dataOption)
		expect := [][]string{
			{"foppa 502", "choice_custom"},
			{"qw.foppa.dk:27502", "qw.foppa.dk:27502"},
		}
		assert.Equal(t, expect, choices)
	})
}
