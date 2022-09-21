package util_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/carmack/internal/pkg/util"
)

func TestContainsAllSubstrings(t *testing.T) {
	assert.False(t, util.ContainsAllSubstrings("qw.foppa.dk:27501", []string{"foo"}))
	assert.False(t, util.ContainsAllSubstrings("qw.foppa.dk:27501", []string{"foppa", "3"}))
	assert.True(t, util.ContainsAllSubstrings("qw.foppa.dk:27501", []string{"foppa"}))
	assert.True(t, util.ContainsAllSubstrings("qw.foppa.dk:27501", []string{"foppa", "2"}))
}

func TestToOptionsChoice(t *testing.T) {
	t.Run("no values", func(t *testing.T) {
		option := []string{}
		expect := &discordgo.ApplicationCommandOptionChoice{Name: "", Value: ""}
		assert.Equal(t, expect, util.ToOptionsChoice(option))
	})

	t.Run("single value", func(t *testing.T) {
		option := []string{"foo"}
		expect := &discordgo.ApplicationCommandOptionChoice{Name: "foo", Value: "foo"}
		assert.Equal(t, expect, util.ToOptionsChoice(option))
	})

	t.Run("two values", func(t *testing.T) {
		option := []string{"foo", "bar"}
		expect := &discordgo.ApplicationCommandOptionChoice{Name: "foo", Value: "bar"}
		assert.Equal(t, expect, util.ToOptionsChoice(option))
	})

	t.Run("more than two values", func(t *testing.T) {
		option := []string{"foo", "bar", "baz"}
		expect := &discordgo.ApplicationCommandOptionChoice{Name: "foo", Value: "bar"}
		assert.Equal(t, expect, util.ToOptionsChoice(option))
	})
}

func TestToOptionChoices(t *testing.T) {
	options := [][]string{
		{"foo"},
		{"foo", "bar"},
		{""},
	}
	expect := []*discordgo.ApplicationCommandOptionChoice{
		{Name: "foo", Value: "foo"},
		{Name: "foo", Value: "bar"},
		{Name: "", Value: ""},
	}
	assert.Equal(t, expect, util.ToOptionChoices(options))
}
