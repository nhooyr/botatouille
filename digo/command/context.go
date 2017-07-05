package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/botatouille/digo/argument"
)

type Context struct {
	DG *discordgo.Session
	S  *argument.Scanner
	M  *discordgo.MessageCreate
}

const (
	normalColor = 0x9933ff
	errorColor  = 0xff3232
)

func (ctx *Context) SendEmbed(e *discordgo.MessageEmbed) (*discordgo.Message, error) {
	// TODO idk about this
	e.Color = normalColor
	return ctx.sendEmbed(e)
}

func (ctx *Context) sendEmbed(e *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.DG.ChannelMessageSendEmbed(ctx.M.ChannelID, e)
}

func (ctx *Context) SendError(err error) (*discordgo.Message, error) {
	return ctx.sendEmbed(&discordgo.MessageEmbed{
		Description: err.Error(),
		Color:       errorColor,
	})
}
