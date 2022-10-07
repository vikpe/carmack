package embed

import (
	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

func FromServer(server qserver.GenericServer) *discordgo.MessageEmbed {
	switch server.Version.GetType() {
	case mvdsv.Name:
		return FromMvdsvServer(convert.ToMvdsv(server))
	case qtv.Name:
		return FromQtvServer(convert.ToQtv(server))
	case qwfwd.Name:
		return FromQwfwdServer(convert.ToQwfwd(server))
	default:
		return FromGenericServer(server)
	}
}
