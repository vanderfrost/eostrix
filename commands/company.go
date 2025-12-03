package commands

import (
	"eostrix/leetcode"
	"eostrix/utils"
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
		utils.ResponseError(s, i, "Invalid difficulty provided.")
		return
	}

	problems, ok := leetcode.ProblemsByCompany[company]
	if !ok || len(problems) == 0 {
		utils.ResponseError(s, i, fmt.Sprintf("No problems found for company %s", company))
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
		utils.ResponseError(s, i, "No matching problems for that difficulty.")
		return
	}

	var sb strings.Builder

	for _, p := range filtered {
		sb.WriteString(fmt.Sprintf("• %s (%s) (%s Frequency)\n%s\n\n",
			p.Title, p.Difficulty, p.Frequency, p.Link))
	}

	utils.Response(s, i, fmt.Sprintf("**%s** %s problems:\n\n", company, difficulty), sb.String())
}
