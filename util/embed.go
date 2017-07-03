package util

import "github.com/bwmarrin/discordgo"

const embedColor = 0x9933ff

func SendEmbed(s *discordgo.Session, channelID string, e *discordgo.MessageEmbed) (*discordgo.Message, error) {
	e.Color = embedColor
	return s.ChannelMessageSendEmbed(channelID, e)
}
