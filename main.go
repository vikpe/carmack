package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

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
		if handler, ok := b.commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
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
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "address",
				Description: "Server address",
				Required:    true,
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
				Content: responseContent,
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