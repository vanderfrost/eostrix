package commands

import (
	"eostrix/leetcode"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleCompanyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	opts := data.Options

	company := strings.ToLower(opts[0].StringValue())
	difficulty := strings.ToLower(opts[1].StringValue())

	switch difficulty {
	case "easy", "medium", "hard", "all":
	default:
		fmt.Println("Invalid difficulty provided.")
		return
	}

	problems, ok := leetcode.ProblemsByCompany[company]
	if !ok || len(problems) == 0 {
		fmt.Printf("No problems found for company %s\n", company)
		return
	}

	// filter by difficulty
	var filtered []*leetcode.Problem
	for _, p := range problems {
		if difficulty == "all" || strings.EqualFold(p.Difficulty, difficulty) {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == 0 {
		fmt.Println("No matching problems for that difficulty.")
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s** %s problems:\n\n",
		company,
		difficulty,
	))

	for _, p := range filtered {
		sb.WriteString(fmt.Sprintf("• [%s] (%s) (%s)\n%s\n\n",
			p.Title, p.Difficulty, p.Frequency, p.Link))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})
}
