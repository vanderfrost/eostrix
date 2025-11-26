package leetcode

import (
	"encoding/json"
	"eostrix/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type leetcodeResponse struct {
	Data struct {
		ActiveDailyCodingChallengeQuestion struct {
			Date     string `json:"date"`
			Link     string `json:"link"`
			Question struct {
				Title      string `json:"title"`
				TitleSlug  string `json:"titleSlug"`
				Difficulty string `json:"difficulty"`
			} `json:"question"`
		} `json:"activeDailyCodingChallengeQuestion"`
	} `json:"data"`
}

type LeetcodeChallenge struct {
	Date       string
	Link       string
	Title      string
	TitleSlug  string
	Difficulty string
}

func parseLeetCodeResponse(resp *leetcodeResponse) LeetcodeChallenge {
	q := resp.Data.ActiveDailyCodingChallengeQuestion
	return LeetcodeChallenge{
		Date:       q.Date,
		Link:       q.Link,
		Title:      q.Question.Title,
		TitleSlug:  q.Question.TitleSlug,
		Difficulty: q.Question.Difficulty,
	}
}

func getDailyChallenge() (LeetcodeChallenge, error) {
	url := "https://leetcode.com/graphql"
	query := `{"query":"{ activeDailyCodingChallengeQuestion { date link question { title titleSlug difficulty } } }"}`

	resp, err := http.Post(url, "application/json", strings.NewReader(query))
	if err != nil {
		return LeetcodeChallenge{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LeetcodeChallenge{}, err
	}

	var lcr leetcodeResponse
	if err := json.Unmarshal(body, &lcr); err != nil {
		return LeetcodeChallenge{}, err
	}

	return parseLeetCodeResponse(&lcr), nil
}

func PostDailyChallenge(session *discordgo.Session) {
	var builder strings.Builder
	cfg := config.ParseConfig()

	challenge, err := getDailyChallenge()
	if err != nil {
		log.Printf("Error fetching challenge: %v", err)
	}

	builder.WriteString(fmt.Sprintf("Challenge Name: %s\n", challenge.Title))
	builder.WriteString(fmt.Sprintf("Date: %s\n", challenge.Date))
	builder.WriteString(fmt.Sprintf("Difficulty: %s\n", challenge.Difficulty))
	builder.WriteString(fmt.Sprintf("Link: \nhttps://leetcode.com%s\n", challenge.Link))

	ping := fmt.Sprintf("<@&%s> ", cfg.LeetcodeRoleID)

	_, err = session.ChannelMessageSendComplex(cfg.DefaultChannel, &discordgo.MessageSend{
		Content: ping,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Daily LeetCode Challenge",
				Description: builder.String(),
				Color:       0xe8a726,
			},
		},
	})
	if err != nil {
		log.Printf("Error sending challenge embed: %v", err)
	}
}
