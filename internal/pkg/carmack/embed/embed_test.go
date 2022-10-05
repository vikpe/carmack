package embed_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/qw-hub-api/types"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qtime"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestFromMvdsvServer(t *testing.T) {
	server := mvdsv.Mvdsv{
		Address: "quake.se:28501",
		Mode:    qmode.Mode("4on4"),
		Title:   "4on4: red (XantoM) [dm2]",
		Status: qstatus.Status{
			Name:        "Started",
			Description: "3 min left",
		},
		Time: qtime.Time{
			Elapsed:   7,
			Total:     10,
			Remaining: 3,
		},
		Players: []qclient.Client{{
			Name:   qstring.New("XantoM"),
			Team:   qstring.New("red"),
			Skin:   "",
			Colors: [2]uint8{13, 13},
			Frags:  2,
			Ping:   38,
			Time:   4,
		},
		},
		PlayerSlots: slots.Slots{
			Used:  1,
			Total: 8,
			Free:  7,
		},
		SpectatorNames: []string{"[ServeMe]"},
		SpectatorSlots: slots.Slots{
			Used:  1,
			Total: 4,
			Free:  3,
		},
		Settings: qsettings.Settings{"map": "dm2", "*gamedir": "qw", "status": "3 min left", "timelimit": "10", "maxclients": "8", "maxspectators": "4", "teamplay": "2", "*version": "MVDSV 0.34", "ktxver": "1.40", "*admin": "spam@foppa.dk"},
		QtvStream: qtvstream.QtvStream{
			Title:          "QUAKE.SE KTX Qtv (1)",
			Url:            "1@quake.se:28000",
			ID:             1,
			Address:        "quake.se:28000",
			SpectatorNames: []string{},
			SpectatorCount: 0,
		},
		Geo: geo.Location{
			CC:          "SE",
			Country:     "Sweden",
			Region:      "Europe",
			City:        "Hagersten",
			Coordinates: [2]float32{59.2885, 17.9612},
		},
		Score: 5,
	}

	expect := &discordgo.MessageEmbed{
		Title:       ":flag_se: quake.se:28501",
		Description: "4on4 on dm2 (started - 3 min left)",
		Color:       0x0c2aac,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Players (1/8)",
				Value: "XantoM",
			},
			{
				Name:  "Spectators (1/4)",
				Value: "[ServeMe]",
			},
			{
				Name:  "QTV",
				Value: "`1@quake.se:28000`",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Admin: spam@foppa.dk [MVDSV 0.34, KTX 1.40]",
		},
	}

	assert.Equal(t, expect, embed.FromMvdsvServer(server))
}

func TestFromStream(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		stream := types.TwitchStream{
			Channel:       "QuakeWorld",
			Url:           "https://twitch.tv/Quakeworld",
			Title:         "1on1: dough vs grl [ztndm3]",
			ViewerCount:   5,
			ServerAddress: "qw.foppa.dk:27502",
		}

		expect := &discordgo.MessageEmbed{
			Title:       "QuakeWorld",
			URL:         "https://twitch.tv/Quakeworld",
			Description: "1on1: dough vs grl [ztndm3]",
			Color:       0xa970ff,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Server",
					Value:  "`qw.foppa.dk:27502`",
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "5 viewers",
			},
		}

		assert.Equal(t, expect, embed.FromStream(stream))
	})

	t.Run("no server address", func(t *testing.T) {
		stream := types.TwitchStream{ServerAddress: ""}
		assert.Empty(t, embed.FromStream(stream).Fields)
	})
}
