package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/log"
	"os/signal"
	"os"
	"syscall"
	"math/rand"
	"github.com/nhooyr/botatouille/util"
	"github.com/nhooyr/botatouille/parser"
)

const (
	token = "MzI2Njk3ODYyNTEzODE5NjQ4.DDtpBQ.j1gzyic1VG3If3yYpUPoWmpVOLw"
	game  = "yo mom"
)

func main() {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
		return
	}

	s.AddHandler(ready)
	s.LogLevel = discordgo.LogWarning
	s.AddHandler(messageCreate)

	err = s.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Print("bot started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-c
	log.Print("exiting")
	err = s.Close()
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmdLine, ok := parser.IsCommand(m.Content)
	if !ok {
		return
	}
	cmd, cmdLine := parser.NextCommand(cmdLine)
	switch cmd {
	case "flip":
		flip(s, m, nil)
	}
}

func flip(s *discordgo.Session, m *discordgo.MessageCreate, _args []string) {
	i := rand.Intn(2)

	e := &discordgo.MessageEmbed{
		Description: coinSides[i],
	}
	util.SendEmbed(s, m.ChannelID, e)
}
