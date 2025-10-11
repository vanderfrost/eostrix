package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commandHandlers = map[string]func(session *discordgo.Session, message *discordgo.MessageCreate){
	// command list
}

func main() {
	config := ParseConfig()

	disc, err := discordgo.New("Bot " + config.SecurityToken)
	if err != nil {
		log.Fatal(err)
	}

	disc.AddHandler(onMessage)
	PostDailyChallenge(disc)

	//open
	disc.Open()

	fmt.Println("bot has started ...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	//close
	defer disc.Close()
}

func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// make sure we only try and read messages that use the built-in command prefix
	if !strings.HasPrefix(message.Content, "/") {
		return
	}

	fields := strings.Fields(message.Content)
	if len(fields) == 0 {
		return
	}

	command := fields[0]
	handler, exists := commandHandlers[command]
	if exists {
		handler(session, message)
	}
}
