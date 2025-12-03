package utils

import "github.com/bwmarrin/discordgo"

// provides some discord related helper functions
// such as response/messaging creation

// use for non-error messages (native slash commands)
func Response(s *discordgo.Session, i *discordgo.InteractionCreate, title, description string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Description: description,
					Color:       0xFFA116,
				},
			},
		},
	})
}

// use for all error messages (native slash commands)
func ResponseError(s *discordgo.Session, i *discordgo.InteractionCreate, description string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Error",
					Description: description,
					Color:       0xFFA116,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

// send message with embed
func SendMessageComplex(s *discordgo.Session, channelID, title, description string) {
	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       title,
				Description: description,
				Color:       0xFFA116,
			},
		},
	})
}

// send message with embed and ping
func SendPingMessageComplex(s *discordgo.Session, channelID, title, ping, description string) {
	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: ping,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       title,
				Description: description,
				Color:       0xFFA116,
			},
		},
	})
}
