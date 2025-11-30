package commands

import "github.com/bwmarrin/discordgo"

func RegisterCommands(s *discordgo.Session) []*discordgo.ApplicationCommand {
	commands := []*discordgo.ApplicationCommand{
		{
			// /company <name> <difficulty>
			Name:        "company",
			Description: "Get LeetCode problems by company",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "name",
					Description: "Company name - eg. ByteDance",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "difficulty",
					Description: "difficulty (easy/medium/hard/all)",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Easy", Value: "easy"},
						{Name: "Medium", Value: "medium"},
						{Name: "Hard", Value: "hard"},
						{Name: "All", Value: "all"},
					},
				},
			},
		},
	}

	for _, cmd := range commands {
		_, _ = s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
	}

	return commands
}
