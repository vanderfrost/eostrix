package commands

import (
	"eostrix/leetcode"
	"eostrix/utils"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// provides a random leetcode problem based on difficulty
// /randlc <difficulty>

func HandleRandCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	opts := data.Options

	difficulty := strings.ToLower(opts[0].StringValue())

	switch difficulty {
	case "easy", "medium", "hard", "all":
	default:
		utils.ResponseError(s, i, "Invalid difficulty provided.")
		return
	}

	list := leetcode.ProblemsByDifficulty[difficulty]
	if len(list) == 0 {
		utils.ResponseError(s, i, "No LeetCode problems found for that difficulty.")
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rng.Intn(len(list))
	problem := list[randomIndex]

	var topics []string
	for _, t := range problem.Topics {
		topics = append(topics, t)
	}

	topicString := "None"
	if len(topics) > 0 {
		topicString = strings.Join(topics, ", ")
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("**Challenge Name:** %s\n", problem.Title))
	builder.WriteString(fmt.Sprintf("**Difficulty:** %s\n", problem.Difficulty))
	builder.WriteString(fmt.Sprintf("**Topics:** %s\n", topicString))
	builder.WriteString(fmt.Sprintf("**Link:** \n%s\n", problem.Link))

	utils.Response(s, i, fmt.Sprintf("Random %s LeetCode Problem", problem.Difficulty), builder.String())
}
