package main

import (
	_ "net/http/pprof"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/botatouille/cmd/botatouille/flip"
	"github.com/nhooyr/botatouille/cmd/botatouille/fortune"
	"github.com/nhooyr/botatouille/cmd/botatouille/music"
	"github.com/nhooyr/botatouille/digo/argument"
	"github.com/nhooyr/botatouille/digo/command"
	"github.com/nhooyr/log"
)

const (
	token = "MzI2Njk3ODYyNTEzODE5NjQ4.DDtpBQ.j1gzyic1VG3If3yYpUPoWmpVOLw"
	game  = "yo mom"
)

var r = command.NewRouter("", "")

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
		return
	}

	dg.AddHandler(ready)
	dg.LogLevel = discordgo.LogWarning
	dg.AddHandler(messageCreate)

	r.Append(music.Command)
	r.Append(fortune.Command)
	r.Append(flip.Command)
	r.Append(music.Command)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Print("bot started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
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

func messageCreate(dg *discordgo.Session, m *discordgo.MessageCreate) {
	cmdLine, ok := command.Is(m.Content)
	if !ok {
		return
	}
	s := argument.NewScanner(cmdLine)
	ctx := &command.Context{dg, s, m}
	err := r.Handle(ctx)
	if err != nil {
		ctx.SendError(err)
		return
	}
}
