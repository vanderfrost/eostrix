package commands

import (
	"eostrix/leetcode"
	"eostrix/utils"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// syntax for command is /topics <category> <difficulty>

func HandleTopicsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	opts := data.Options

	topic := strings.ToLower(opts[0].StringValue())
	difficulty := strings.ToLower(opts[1].StringValue())

	problems, ok := leetcode.ProblemsByTopic[topic]
	if !ok || len(problems) == 0 {
		utils.ResponseError(s, i, fmt.Sprintf("No problems found for topic %s", topic))
		return
	}

	// filter by difficulty
	var filtered []*leetcode.Problem
	for _, p := range problems {
		if strings.EqualFold(p.Difficulty, difficulty) {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == 0 {
		utils.ResponseError(s, i, "No matching problems for that difficulty.")
		return
	}

	// envoke a ten problem limit
	// (just to test and see if the command works without pagination or autocomplete)
	limit := 10
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}

	var sb strings.Builder

	for _, p := range filtered {
		sb.WriteString(fmt.Sprintf("• %s (%s) (%s Frequency)\n%s\n\n",
			p.Title, p.Difficulty, p.Frequency, p.Link))
	}

	utils.Response(s, i, fmt.Sprintf("%s Topic Problems", topic), sb.String())
}
