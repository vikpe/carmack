package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

type Bot struct {
	session         *discordgo.Session
	guildID         string
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	OnReady         func(s *discordgo.Session)
}

func NewBot(token string, guildId string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)

	if err != nil {
		return &Bot{}, err
	}

	return &Bot{
		session:         session,
		guildID:         guildId,
		commands:        make([]*discordgo.ApplicationCommand, 0),
		commandHandlers: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), 0),
		OnReady: func(s *discordgo.Session) {
			log.Println(fmt.Sprintf("%s is ready", s.State.User.Username))
		},
	}, nil
}

func (b *Bot) Start() {
	err := b.session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer b.session.Close()

	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		b.OnReady(s)
	})

	b.RegisterCommands()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
	b.UnregisterCommands()
}

func (b *Bot) RegisterCommands() {
	for _, v := range b.commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, b.guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%s' command: %s", v.Name, err)
		}
	}

	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if handler, ok := b.commandHandlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			}

		case discordgo.InteractionApplicationCommandAutocomplete:
			optionMap := toOptionsMap(i.ApplicationCommandData().Options)

			var choices []*discordgo.ApplicationCommandOptionChoice

			if optionMap["address"].Focused {
				inputValue := optionMap["address"].StringValue()

				if inputValue != "" {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  inputValue,
						Value: "choice_custom",
					})
				}

				for _, address := range GetServerAddresses(inputValue) {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name: address, Value: address,
					})
				}
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices,
				},
			})
			if err != nil {
				panic(err)
			}
		}
	})
}

func (b *Bot) UnregisterCommands() {
	log.Println("Removing commands...")

	registeredCommands, err := b.session.ApplicationCommands(b.session.State.User.ID, b.guildID)
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, b.guildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func (b *Bot) AddCommand(cmd *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	b.commands = append(b.commands, cmd)
	b.commandHandlers[cmd.Name] = handler
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	bot, err := NewBot(
		os.Getenv("BOT_TOKEN"),
		os.Getenv("GUILD_ID"),
	)

	if err != nil {
		log.Fatal("unable to create bot", err)
		return
	}

	bot.AddCommand(&discordgo.ApplicationCommand{
		Name:        "server",
		Description: "server command",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "address",
				Description:  "Server address",
				Required:     true,
				Autocomplete: true,
			},
		},
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		optionMap := toOptionsMap(i.ApplicationCommandData().Options)
		server, err := GetServerInfo(optionMap["address"].StringValue())
		responseContent := ""

		if err != nil {
			responseContent = err.Error()
		} else {
			responseContent = fmt.Sprintf("%s - %s", server.Address, server.Title)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: responseContent,
			},
		})
	})

	bot.AddCommand(&discordgo.ApplicationCommand{
		Name:        "find",
		Description: "find player",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Player name",
				Required:    true,
			},
		},
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		optionMap := toOptionsMap(i.ApplicationCommandData().Options)
		playerName := optionMap["name"].StringValue()
		server, err := FindPlayer(playerName)

		if err != nil {
			log.Println("ERROR", err)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				}})

			return
		}

		log.Println("SUCCESS", err)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("%s is playing at %s (%s)", playerName, server.Address, server.Title),
				/*Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "No",
								Style:    discordgo.DangerButton,
								Disabled: false,
								CustomID: "fd_no",
							},
							discordgo.Button{
								Label: "FTP link",
								Style: discordgo.LinkButton,
								URL:   "ftp://foo.bar:21",
							},
						},
					},
				},*/
			},
		})
	})

	bot.Start() // blocking operation
}

func toOptionsMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

func GetServerInfo(address string) (*mvdsv.Mvdsv, error) {
	url := fmt.Sprintf("https://hubapi.quakeworld.nu/v2/servers/%s", address)
	resp, err := resty.New().R().SetResult(&mvdsv.Mvdsv{}).Get(url)

	if err != nil {
		err = errors.New("unable to fetch server information")
	}

	if resp.StatusCode() != http.StatusOK {
		err = errors.New(resp.String())
	}

	if err != nil {
		return &mvdsv.Mvdsv{}, err
	}

	server := resp.Result().(*mvdsv.Mvdsv)

	return server, nil
}

func GetMvdsvServers(queryParams map[string]string) []mvdsv.Mvdsv {
	serversUrl := "https://hubapi.quakeworld.nu/v2/servers/mvdsv"
	resp, err := resty.New().R().SetResult([]mvdsv.Mvdsv{}).SetQueryParams(queryParams).Get(serversUrl)

	if err != nil {
		fmt.Println("server fetch error", err.Error())
		return make([]mvdsv.Mvdsv, 0)
	}

	servers := resp.Result().(*[]mvdsv.Mvdsv)
	return *servers
}

func FindPlayer(pattern string) (mvdsv.Mvdsv, error) {
	const minFindLength = 2

	if len(pattern) < minFindLength {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`provide at least %d characters.`, minFindLength))
	}

	if !strings.Contains(pattern, "*") {
		pattern = fmt.Sprintf("*%s*", pattern)
	}

	servers := GetMvdsvServers(map[string]string{"has_player": pattern})

	if 0 == len(servers) {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`player "%s" not found.`, pattern))
	}

	return servers[0], nil
}

func GetServerAddresses(needle string) []string {
	if "" == needle {
		return make([]string, 0)
	}

	all := []string{
		"clanrot.org:28501", "clanrot.org:28502", "clanrot.org:28503", "clanrts.com:27500", "clanrts.com:27500", "clanrts.com:27501", "clanrts.com:27501", "clanrts.com:27502", "clanrts.com:27503", "cyberdemon.co.uk:28501", "dal.spawnfrag.com:28501", "dal.spawnfrag.com:28502", "dal.spawnfrag.com:28503", "dal.spawnfrag.com:28504", "eskaldar.ru:28501", "ffa.kinsky.io:27500", "fr.predze.dk:28501", "fr.predze.dk:28502", "fr.predze.dk:28503", "ie.predze.dk:28501", "ie.predze.dk:28502", "ie.predze.dk:28503", "ie.qwsrv.com:28501", "ie.qwsrv.com:28502", "in3.antilag.predze.dk:28501", "in3.antilag.predze.dk:28502", "in3.predze.dk:28501", "in3.predze.dk:28502", "in3.predze.dk:28503", "io.qwsrv.com:28501", "io.qwsrv.com:28502", "isf-clan.com:27166", "london.badplace.eu:28501", "london.badplace.eu:28502", "london.badplace.eu:28503", "london.badplace.eu:28504", "neppomuk.ch:27500", "neppomuk.ch:27502", "neppomuk.ch:27505", "nl.aye.wtf:28501", "nl.aye.wtf:28502", "nl.aye.wtf:28503", "nl.aye.wtf:28504", "nl.aye.wtf:28505", "nl.aye.wtf:28506", "nl.aye.wtf:28507", "nl.aye.wtf:28508", "nl.kinsky.io:28501", "nl.kinsky.io:28502", "nl.kinsky.io:28503", "nl.kinsky.io:28504", "nl2.badplace.eu:28501", "nl2.badplace.eu:28502", "nl2.badplace.eu:28503", "nl2.badplace.eu:28504", "nl2.badplace.eu:28505", "nl2.badplace.eu:28506", "nq.karancssag.info:28501", "ow.irc.ax:28501", "ow.irc.ax:28502", "ow.irc.ax:28503", "ow.irc.ax:28504", "pl.dm6.uk:28501", "pl.dm6.uk:28502", "play.quake1.pl:27500", "play.quake1.pl:27501", "play.quake1.pl:27510", "play.quake1.pl:27520", "pummelator.com:28502", "pummelator.com:28503", "pummelator.com:28504", "quake.faked.org:27500", "quakecon.guam:28503", "quakecon.guam:28504", "quakeworld.fi:27500", "quakeworld.fi:28501", "quakeworld.fi:28502", "quakeworld.fi:28503", "quakeworld.fi:28504", "qw.0f.se:28501", "qw.0f.se:28502", "qw.0f.se:28503", "qw.0f.se:28504", "qw.0f.se:28505", "qw.0f.se:28506", "qw.0f.se:28507", "qw.0f.se:28508", "qw.0f.se:28509", "qw.0f.se:28510", "qw.fnu.nu:28501", "qw.foppa.dk:27500", "qw.foppa.dk:27501", "qw.foppa.dk:27502", "qw.foppa.dk:27503", "qw.foppa.dk:27504", "qw.foppa.dk:27505", "qw.irc.ax:28501", "qw.irc.ax:28502", "qw.irc.ax:28503", "qw.irc.ax:28504", "qw.irc.ax:28505", "qw.irc.ax:28506", "qw.irc.ax:28507", "qw.irc.ax:28508", "qw.irc.ax:28509", "qw.irc.ax:28510", "qw.kirril.com:28501", "qw.kirril.com:28502", "qw.kirril.com:28503", "qw2.ru:27500", "qw2.ru:27501", "qw2.ru:27502", "qw2.ru:27503", "sbct.ru:12345", "se.imaginary.zone:28501", "se.imaginary.zone:28502", "se.imaginary.zone:28503", "se.imaginary.zone:28504", "shkn.ws:28501", "shkn.ws:28502", "shkn.ws:28503", "shkn.ws:28504", "soskif.li:28501", "suddendeath.nu:28501", "suddendeath.nu:28502", "suddendeath.nu:28503", "suddendeath.nu:28504", "suddendeath.nu:28505", "suddendeath.nu:28506", "suntzu.cz:27500", "troopers.fi:28001", "troopers.fi:28002", "troopers.fi:28003", "tummen.se:27500", "wa.b1aze.com:28501", "wa.b1aze.com:28502", "ziewa.pl:28501", "ziewa.pl:28502", "ziewa.pl:28503", "ziewa.pl:28504",
	}

	var result []string
	words := strings.Split(needle, " ")
	const maxResults = 10

	for _, address := range all {
		if !containsAllSubstrings(address, words) {
			continue
		}

		result = append(result, address)

		if len(result) == maxResults {
			break
		}
	}

	return result
}

func containsAllSubstrings(haystack string, needles []string) bool {
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			return false
		}
	}

	return true
}
