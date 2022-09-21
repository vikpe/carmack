package autocomplete

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/util"
)

func ServerAddress(option *discordgo.ApplicationCommandInteractionDataOption) [][]string {
	choices := make([][]string, 0)
	inputValue := option.StringValue()

	if inputValue != "" {
		choices = append(choices, []string{inputValue, "choice_custom"})
	}

	for _, address := range getServerAddresses(inputValue) {
		choices = append(choices, []string{address, address})
	}

	return choices
}

func getServerAddresses(needle string) []string {
	if "" == needle {
		return make([]string, 0)
	}

	all := []string{
		"clanrot.org:28501", "clanrot.org:28502", "clanrot.org:28503", "clanrts.com:27500", "clanrts.com:27500", "clanrts.com:27501", "clanrts.com:27501", "clanrts.com:27502", "clanrts.com:27503", "cyberdemon.co.uk:28501", "dal.spawnfrag.com:28501", "dal.spawnfrag.com:28502", "dal.spawnfrag.com:28503", "dal.spawnfrag.com:28504", "eskaldar.ru:28501", "ffa.kinsky.io:27500", "fr.predze.dk:28501", "fr.predze.dk:28502", "fr.predze.dk:28503", "ie.predze.dk:28501", "ie.predze.dk:28502", "ie.predze.dk:28503", "ie.qwsrv.com:28501", "ie.qwsrv.com:28502", "in3.antilag.predze.dk:28501", "in3.antilag.predze.dk:28502", "in3.predze.dk:28501", "in3.predze.dk:28502", "in3.predze.dk:28503", "io.qwsrv.com:28501", "io.qwsrv.com:28502", "isf-clan.com:27166", "london.badplace.eu:28501", "london.badplace.eu:28502", "london.badplace.eu:28503", "london.badplace.eu:28504", "neppomuk.ch:27500", "neppomuk.ch:27502", "neppomuk.ch:27505", "nl.aye.wtf:28501", "nl.aye.wtf:28502", "nl.aye.wtf:28503", "nl.aye.wtf:28504", "nl.aye.wtf:28505", "nl.aye.wtf:28506", "nl.aye.wtf:28507", "nl.aye.wtf:28508", "nl.kinsky.io:28501", "nl.kinsky.io:28502", "nl.kinsky.io:28503", "nl.kinsky.io:28504", "nl2.badplace.eu:28501", "nl2.badplace.eu:28502", "nl2.badplace.eu:28503", "nl2.badplace.eu:28504", "nl2.badplace.eu:28505", "nl2.badplace.eu:28506", "nq.karancssag.info:28501", "ow.irc.ax:28501", "ow.irc.ax:28502", "ow.irc.ax:28503", "ow.irc.ax:28504", "pl.dm6.uk:28501", "pl.dm6.uk:28502", "play.quake1.pl:27500", "play.quake1.pl:27501", "play.quake1.pl:27510", "play.quake1.pl:27520", "pummelator.com:28502", "pummelator.com:28503", "pummelator.com:28504", "quake.faked.org:27500", "quakecon.guam:28503", "quakecon.guam:28504", "quakeworld.fi:27500", "quakeworld.fi:28501", "quakeworld.fi:28502", "quakeworld.fi:28503", "quakeworld.fi:28504", "qw.0f.se:28501", "qw.0f.se:28502", "qw.0f.se:28503", "qw.0f.se:28504", "qw.0f.se:28505", "qw.0f.se:28506", "qw.0f.se:28507", "qw.0f.se:28508", "qw.0f.se:28509", "qw.0f.se:28510", "qw.fnu.nu:28501", "qw.foppa.dk:27500", "qw.foppa.dk:27501", "qw.foppa.dk:27502", "qw.foppa.dk:27503", "qw.foppa.dk:27504", "qw.foppa.dk:27505", "qw.irc.ax:28501", "qw.irc.ax:28502", "qw.irc.ax:28503", "qw.irc.ax:28504", "qw.irc.ax:28505", "qw.irc.ax:28506", "qw.irc.ax:28507", "qw.irc.ax:28508", "qw.irc.ax:28509", "qw.irc.ax:28510", "qw.kirril.com:28501", "qw.kirril.com:28502", "qw.kirril.com:28503", "qw2.ru:27500", "qw2.ru:27501", "qw2.ru:27502", "qw2.ru:27503", "sbct.ru:12345", "se.imaginary.zone:28501", "se.imaginary.zone:28502", "se.imaginary.zone:28503", "se.imaginary.zone:28504", "shkn.ws:28501", "shkn.ws:28502", "shkn.ws:28503", "shkn.ws:28504", "soskif.li:28501", "suddendeath.nu:28501", "suddendeath.nu:28502", "suddendeath.nu:28503", "suddendeath.nu:28504", "suddendeath.nu:28505", "suddendeath.nu:28506", "suntzu.cz:27500", "troopers.fi:28001", "troopers.fi:28002", "troopers.fi:28003", "tummen.se:27500", "wa.b1aze.com:28501", "wa.b1aze.com:28502", "ziewa.pl:28501", "ziewa.pl:28502", "ziewa.pl:28503", "ziewa.pl:28504",
	}

	var result []string
	words := strings.Split(needle, " ")
	const maxResults = 20

	for _, address := range all {
		if !util.ContainsAllSubstrings(address, words) {
			continue
		}

		result = append(result, address)

		if len(result) == maxResults {
			break
		}
	}

	return result
}
