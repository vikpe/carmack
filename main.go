package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
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
		Name:        "ping",
		Description: "ping command",
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there! Congratulations, you just executed your first slash command",
			},
		})
	})

	bot.Start() // blocking operation
}
