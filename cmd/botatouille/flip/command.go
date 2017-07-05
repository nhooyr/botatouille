package flip

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/botatouille/digo/command"
)

var Command = command.NewCommand("flip", "", "bombi", flip)

var coinSides = []string{"heads", "tails"}

func flip(ctx *command.Context) error  {
	i := rand.Intn(2)
	e := &discordgo.MessageEmbed{
		Description: coinSides[i],
	}
	_, err := ctx.SendEmbed(e)
	return err
}