package discordbot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/vikpe/carmack/internal/pkg/util"
)

type AutocompleteHandler func(option *discordgo.ApplicationCommandInteractionDataOption) [][]string

type CommandHandler func(i *discordgo.InteractionCreate) *discordgo.InteractionResponse

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
	for _, command := range b.commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, b.guildID, command)
		if err != nil {
			log.Panicf("Cannot create '%s' command: %s", command.Name, err)
		}
	}

	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			log.Println(fmt.Sprintf("## command: %s", i.ApplicationCommandData().Name))

			if handler, ok := b.commandHandlers[i.ApplicationCommandData().Name]; ok {
				interactionResponse := handler(i)
				err := s.InteractionRespond(i.Interaction, interactionResponse)
				if err != nil {
					log.Fatal(err)
					return
				}
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
						Choices: util.ToOptionChoices(choices),
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

func (b *Bot) AddCommand(command *discordgo.ApplicationCommand, handler CommandHandler) {
	b.commands = append(b.commands, command)
	b.commandHandlers[command.Name] = handler
}

func (b *Bot) AddAutocompleteHandler(key string, handler AutocompleteHandler) {
	b.autocompleteHandlers[key] = handler
}
