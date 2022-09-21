package autocomplete_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/assert/v2"
	"github.com/vikpe/carmack/internal/pkg/carmack/autocomplete"
)

func TestServerAddress(t *testing.T) {
	dataOption := &discordgo.ApplicationCommandInteractionDataOption{
		Name:    "address",
		Type:    discordgo.ApplicationCommandOptionString,
		Value:   "foppa",
		Focused: true,
	}

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
}
