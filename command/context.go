package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/botatouille/argument"
)

type Context struct {
	DG *discordgo.Session
	S  *argument.Scanner
	M  *discordgo.MessageCreate
}