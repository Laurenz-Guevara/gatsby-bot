package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutting down or not")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	// Load environment variables from .env file
	env := godotenv.Load()
	if env != nil {
		log.Fatalf("Error loading .env file: %v", env)
	}

	// Set BotToken from environment variable
	botToken := os.Getenv("BOT_TOKEN")

	// Initialize Discord session
	var initErr error
	s, initErr = discordgo.New("Bot " + botToken)
	if initErr != nil {
		log.Fatalf("Invalid bot parameters: %v", initErr)
	}
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name: "mood",
			Description: "Work out if Gatsby is happy, sad or something else.",
		},
		{
			Name:        "pet",
			Description: "Subcommands and command groups example",
			Options: []*discordgo.ApplicationCommandOption{
				// {
				// 	Name:        "subcommand-group",
				// 	Description: "Subcommands group",
				// 	Options: []*discordgo.ApplicationCommandOption{
				// 		{
				// 			Name:        "nested-subcommand",
				// 			Description: "Nested subcommand",
				// 			Type:        discordgo.ApplicationCommandOptionSubCommand,
				// 		},
				// 	},
				// 	Type: discordgo.ApplicationCommandOptionSubCommandGroup,
				// },
				{
					Name:        "gatsby",
					Description: "Give gatsby a head pat",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mood": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Gatsby's mood is cheerful!",
				},
			})
		},
		"pet": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			content := ""

			// As you can see, names of subcommands (nested, top-level)
			// and subcommand groups are provided through the arguments.
			switch options[0].Name {
			case "gatsby":
				content = "meow *You petted Gatsby! +1 Rep*"
				
			// case "subcommand-group":
			// 	options = options[0].Options
			// 	switch options[0].Name {
			// 	case "nested-subcommand":
			// 		content = "Nice, now you know how to execute nested commands too"
			// 	default:
			// 		content = "Oops, something went wrong.\n" +
			// 			"Hol' up, you aren't supposed to see this message."
			// 	}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
