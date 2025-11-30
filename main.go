package main

import (
	"eostrix/commands"
	"eostrix/config"
	"eostrix/leetcode"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config := config.ParseConfig()

	disc, err := discordgo.New("Bot " + config.SecurityToken)
	if err != nil {
		log.Fatal(err)
	}

	//open
	disc.Open()

	//register commands
	commands.RegisterCommands(disc)
	disc.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			switch i.ApplicationCommandData().Name {
			case "company":
				commands.HandleCompanyCommand(s, i)
			}
		}
	})

	//schedule daily post
	leetcode.ScheduleMidnightUTCEvent(func() {
		leetcode.PostDailyChallenge(disc)
	})

	//load company LC data
	_, err = leetcode.LoadAllProblems("data")
	if err != nil {
		log.Fatal("failed to load company problems:", err)
	}

	fmt.Println("bot has started ...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	//close
	defer disc.Close()
}
