package main

import (
	"eostrix/commands"
	"eostrix/config"
	"eostrix/leetcode"
	"eostrix/utils"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config := config.ParseConfig()

	disc, err := discordgo.New("Bot " + config.SecurityToken)
	if err != nil {
		log.Fatal(err)
	}

	if err := disc.Open(); err != nil {
		log.Fatal(err)
	}
	defer disc.Close()

	initHandlers(disc)

	if err := loadFeatures(disc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("bot has started ...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func initHandlers(disc *discordgo.Session) {
	commands.RegisterCommands(disc)
	disc.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			switch i.ApplicationCommandData().Name {
			case "company":
				commands.HandleCompanyCommand(s, i)
			case "randlc":
				commands.HandleRandCommand(s, i)
			case "topics":
				commands.HandleTopicsCommand(s, i)
			}
		}
		if i.Type == discordgo.InteractionMessageComponent {
			cid := i.MessageComponentData().CustomID
			if strings.HasPrefix(cid, "company_next") {
				commands.HandleCompanyPageChange(s, i, +1)
			} else if strings.HasPrefix(cid, "company_prev") {
				commands.HandleCompanyPageChange(s, i, -1)
			}
			return
		}
	})

	disc.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommandAutocomplete {

			switch i.ApplicationCommandData().Name {
			case "company":
				commands.CompanyAutocomplete(s, i)
			case "topics":
				commands.TopicsAutocomplete(s, i)
			}
		}
	})
}

func loadFeatures(disc *discordgo.Session) error {
	utils.ScheduleMidnightUTCEvent(func() {
		leetcode.PostDailyChallenge(disc)
	})

	_, err := leetcode.LoadAllProblems("data")
	return err
}
