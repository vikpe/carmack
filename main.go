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
				customValue := optionMap["address"].StringValue()
				choices = []*discordgo.ApplicationCommandOptionChoice{
					{Name: "qw.foppa.dk:27501", Value: "qw.foppa.dk:27501"},
					{Name: "qw.foppa.dk:27502", Value: "qw.foppa.dk:27502"},
					{Name: "qw.foppa.dk:27503", Value: "qw.foppa.dk:27503"},
					{Name: "qw.foppa.dk:27504", Value: "qw.foppa.dk:27504"},
					{Name: "qw.foppa.dk:27505", Value: "qw.foppa.dk:27505"},
				}
				if customValue != "" {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  customValue,
						Value: "choice_custom",
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
