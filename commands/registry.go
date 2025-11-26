package commands

import "github.com/bwmarrin/discordgo"

func RegisterCommands(s *discordgo.Session) []*discordgo.ApplicationCommand {
	commands := []*discordgo.ApplicationCommand{
		{
			// example for now
			Name:        "company",
			Description: "Get LeetCode problems by company",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "name",
					Description: "Company name - eg. ByteDance",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}

	for _, cmd := range commands {
		_, _ = s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
	}

	return commands
}
