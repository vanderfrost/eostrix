package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

func GetDailyChallenge() (LeetcodeChallenge, error) {
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
