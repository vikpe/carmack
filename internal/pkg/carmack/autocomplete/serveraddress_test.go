package autocomplete_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
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
	t.Run("empty needle", func(t *testing.T) {
		dataOption := getDataOption("")
		choices := autocomplete.ServerAddress(dataOption)
		expect := [][]string{}
		assert.Equal(t, expect, choices)
	})

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

	t.Run("more than 20 results", func(t *testing.T) {
		dataOption := getDataOption(":")
		choices := autocomplete.ServerAddress(dataOption)
		expect := [][]string{
			{":", "choice_custom"},
			{"clanrot.org:28501", "clanrot.org:28501"},
			{"clanrot.org:28502", "clanrot.org:28502"},
			{"clanrot.org:28503", "clanrot.org:28503"},
			{"clanrts.com:27500", "clanrts.com:27500"},
			{"clanrts.com:27500", "clanrts.com:27500"},
			{"clanrts.com:27501", "clanrts.com:27501"},
			{"clanrts.com:27501", "clanrts.com:27501"},
			{"clanrts.com:27502", "clanrts.com:27502"},
			{"clanrts.com:27503", "clanrts.com:27503"},
			{"cyberdemon.co.uk:28501", "cyberdemon.co.uk:28501"},
			{"dal.spawnfrag.com:28501", "dal.spawnfrag.com:28501"},
			{"dal.spawnfrag.com:28502", "dal.spawnfrag.com:28502"},
			{"dal.spawnfrag.com:28503", "dal.spawnfrag.com:28503"},
			{"dal.spawnfrag.com:28504", "dal.spawnfrag.com:28504"},
			{"eskaldar.ru:28501", "eskaldar.ru:28501"},
			{"ffa.kinsky.io:27500", "ffa.kinsky.io:27500"},
			{"fr.predze.dk:28501", "fr.predze.dk:28501"},
			{"fr.predze.dk:28502", "fr.predze.dk:28502"},
			{"fr.predze.dk:28503", "fr.predze.dk:28503"},
			{"ie.predze.dk:28501", "ie.predze.dk:28501"},
		}
		assert.Equal(t, expect, choices)
	})
}
