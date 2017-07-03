package main

import (
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/botatouille/argument"
	"github.com/nhooyr/botatouille/command"
	"github.com/nhooyr/botatouille/util"
	"github.com/nhooyr/log"
)

const (
	token = "MzI2Njk3ODYyNTEzODE5NjQ4.DDtpBQ.j1gzyic1VG3If3yYpUPoWmpVOLw"
	game  = "yo mom"
)

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
		return
	}

	dg.AddHandler(ready)
	dg.LogLevel = discordgo.LogWarning
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Print("bot started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-c
	log.Print("exiting")
	err = dg.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func ready(s *discordgo.Session, _r *discordgo.Ready) {
	err := s.UpdateStatus(0, game)
	if err != nil {
		log.Fatal(err)
	}
}

var coinSides = []string{"heads", "tails"}

func messageCreate(dg *discordgo.Session, m *discordgo.MessageCreate) {
	cmdLine, ok := command.Is(m.Content)
	if !ok {
		return
	}
	s := argument.NewScanner(cmdLine)
	err := s.Scan()
	switch err {
	case argument.ErrUnexpectedBackspace:
		// TODO link to some docs explaining how answers work
		dg.ChannelMessageSend(m.ChannelID, "You need to have a character after the backspace.")
		return
	case io.EOF:
		dg.ChannelMessageSend(m.ChannelID, "You must give me a command name.")
		return
	}
	ctx := &command.Context{dg, s, m}
	switch s.Token() {
	case "flip":
		flip(ctx)
	case "fortune":
		fortune(ctx)
	case "m":
		music(ctx)
	case "c":
		custom_command(ctx)
	}
}

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

func fortune(ctx *command.Context) {
	i := rand.Intn(len(answers))
	question := ctx.S.Rest()
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
	util.SendEmbed(ctx.DG, ctx.M.ChannelID, e)
}

func flip(ctx *command.Context) {
	i := rand.Intn(2)
	e := &discordgo.MessageEmbed{
		Description: coinSides[i],
	}
	util.SendEmbed(ctx.DG, ctx.M.ChannelID, e)
}