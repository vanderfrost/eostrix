package commands

import (
	"eostrix/leetcode"
	"eostrix/utils"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HandleCompanyCommand sends the first page of company problems
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

	pageData := &leetcode.PageData{
		Problems: filtered,
		Page:     0,
		PageSize: 10,
	}
	leetcode.StorePage(i.Member.User.ID, pageData)

	renderCompanyPage(s, i, pageData, company, difficulty, true)
}

func CompanyAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	userInput := strings.ToLower(data.Options[0].StringValue())
	var suggestions []*discordgo.ApplicationCommandOptionChoice

	for _, vc := range leetcode.ValidCompanies {
		if strings.Contains(strings.ToLower(vc), userInput) {
			suggestions = append(suggestions,
				&discordgo.ApplicationCommandOptionChoice{
					Name:  vc,
					Value: vc,
				})
		}
		if len(suggestions) == 25 {
			break
		}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: suggestions,
		},
	})
}

// renderCompanyPage handles both first page posting and editing of all other pages
func renderCompanyPage(s *discordgo.Session, i *discordgo.InteractionCreate, data *leetcode.PageData, company, difficulty string, first bool) {
	start := data.Page * data.PageSize
	end := min(start+data.PageSize, len(data.Problems))

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"**Company: %s**\n**Difficulty: %s**\n**Page %d / %d**\n\n",
		company,
		difficulty,
		data.Page+1,
		(len(data.Problems)+data.PageSize-1)/data.PageSize,
	))

	for _, p := range data.Problems[start:end] {
		sb.WriteString(fmt.Sprintf("• %s (%s) %s freq\n%s\n\n",
			p.Title, p.Difficulty, p.Frequency, p.Link))
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					CustomID: fmt.Sprintf("company_prev|%s|%s", company, difficulty),
					Label:    "Prev",
					Style:    discordgo.PrimaryButton,
					Disabled: data.Page == 0,
				},
				&discordgo.Button{
					CustomID: fmt.Sprintf("company_next|%s|%s", company, difficulty),
					Label:    "Next",
					Style:    discordgo.PrimaryButton,
					Disabled: (data.Page+1)*data.PageSize >= len(data.Problems),
				},
			},
		},
	}

	if first {
		utils.ResponseComponents(s, i, sb.String(), components)
	} else {
		utils.ResponseComponentsEdit(s, i, sb.String(), components)
	}
}

func HandleCompanyPageChange(s *discordgo.Session, i *discordgo.InteractionCreate, delta int) {
	if i.Member == nil || i.Member.User == nil {
		utils.ResponseError(s, i, "cannot retrieve user info.")
		return
	}

	user := i.Member.User.ID
	data, ok := leetcode.GetPage(user)
	if !ok {
		utils.ResponseError(s, i, "pagination expired or missing.")
		return
	}

	parts := strings.Split(i.MessageComponentData().CustomID, "|")
	if len(parts) != 3 {
		utils.ResponseError(s, i, "invalid button data.")
		return
	}
	company := parts[1]
	difficulty := parts[2]

	data.Page += delta
	if data.Page < 0 {
		data.Page = 0
	}
	if data.Page*data.PageSize >= len(data.Problems) {
		data.Page -= delta
	}

	renderCompanyPage(s, i, data, company, difficulty, false)
}
