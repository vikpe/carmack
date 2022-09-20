package discordbot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/util"
)

type AutocompleteHandler func(option *discordgo.ApplicationCommandInteractionDataOption) []*discordgo.ApplicationCommandOptionChoice

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Bot struct {
	session              *discordgo.Session
	guildID              string
	commands             []*discordgo.ApplicationCommand
	commandHandlers      map[string]CommandHandler
	autocompleteHandlers map[string]AutocompleteHandler
	OnReady              func(s *discordgo.Session)
}

func New(token string, guildId string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)

	if err != nil {
		return &Bot{}, err
	}

	return &Bot{
		session:              session,
		guildID:              guildId,
		commands:             make([]*discordgo.ApplicationCommand, 0),
		commandHandlers:      make(map[string]CommandHandler, 0),
		autocompleteHandlers: make(map[string]AutocompleteHandler, 0),
		OnReady: func(s *discordgo.Session) {
			log.Println(fmt.Sprintf("%s is ready", s.State.User.Username))
		},
	}, nil
}

func (b *Bot) Start() {
	log.Println("Start()")

	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		b.OnReady(s)
	})

	err := b.session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer b.session.Close()

	b.RegisterCommands()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
	b.UnregisterCommands()
}

func (b *Bot) RegisterCommands() {
	log.Println("RegisterCommands()")
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
			optionMap := util.ToOptionsMap(i.ApplicationCommandData().Options)
			commandDataName := "undefined"

			for k := range optionMap {
				if optionMap[k].Focused {
					commandDataName = k
				}
			}

			if handler, ok := b.autocompleteHandlers[commandDataName]; ok {
				choices := handler(optionMap[commandDataName])

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionApplicationCommandAutocompleteResult,
					Data: &discordgo.InteractionResponseData{
						Choices: choices,
					},
				})
				if err != nil {
					log.Fatal(err)
				}

			} else {
				log.Println(fmt.Sprintf(`no autocomplete handler defined for "%s"`, commandDataName))
			}

		}
	})
}

func (b *Bot) UnregisterCommands() {
	log.Println("UnregisterCommands()")

	registeredCommands, err := b.session.ApplicationCommands(b.session.State.User.ID, b.guildID)
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		log.Println(fmt.Sprintf(`removing command "%s"`, v.Name))

		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, b.guildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func (b *Bot) AddCommand(command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	b.commands = append(b.commands, command)
	b.commandHandlers[command.Name] = handler
}

func (b *Bot) AddAutocompleteHandler(key string, handler AutocompleteHandler) {
	b.autocompleteHandlers[key] = handler
}
