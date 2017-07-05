package fortune

import (
	"strings"

	"math/rand"

	"github.com/bwmarrin/discordgo"
	"errors"
	"github.com/nhooyr/botatouille/digo/command"
)

var answers = []string{
	"It's quite possible.",
	"Only if Aaron's dad is a pimp.",
	"Only If Chris Wingy gave up his right testicle.",
	"Only the flying spaghetti monster knows.",
	"Of course.",
	"Absolutely not.",
	"Sacrifice yourself to Kyle's cows and the answer will be before you.",
	"Test it.",
	"Great question, unfortunately I have no idea.",
	"Pray before bambi's cows and u will find the answer.",
	"Sock my dog lock.",
	"Prepare yourself to be upon the gods for only they know.",
	"Careful, this is a dangerous question.",
	"Do you take me for a fool?",
	"Do you think I'm god...?",
	"Kinky question, almost as kinky as kiky",
}

var Command = command.NewCommand("fortune", "{question}", "boo", fortune)

func fortune(ctx *command.Context) error {
	i := rand.Intn(len(answers))
	question := ctx.S.Rest()
	if question == "" {
		return errors.New("What is the question you seek?")
	}
	if !strings.HasSuffix(question, "?") {
		question = question + "?"
	}
	e := &discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Question",
				Value: question,
			},
			{
				Name:  "Answer",
				Value: answers[i],
			},
		},
	}
	_, err := ctx.SendEmbed(e)
	return err
}
