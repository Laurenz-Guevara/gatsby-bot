package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const prefix string = "!c"

func main() {
	fmt.Println("Starting Gatsby Cat Bot")
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	commandResponses := map[string]string{
		"Cat":  "meow",
		"Dog":  "bark",
		"bird": "tweet",
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		//Map
		if response, exists := commandResponses[m.Content]; exists {
			s.ChannelMessageSend(m.ChannelID, response)
		} 

		args := strings.Split(m.Content, " ")

		if args[0] != prefix {
			return
		}

		input := strings.ToLower(args[1])

		if input == "commands" {
			for key := range commandResponses {
				s.ChannelMessageSend(m.ChannelID, key)
			}
		}

		if input == "owner" {
			s.ChannelMessageSend(m.ChannelID, "Zushi and Cookiee are my owners.")
		}
	})
 
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println("The bot is online!")
	
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}