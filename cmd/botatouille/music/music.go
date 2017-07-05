package music

import (
	"errors"

	"os/exec"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/botatouille/digo/argument"
	"github.com/nhooyr/botatouille/digo/command"
	"github.com/nhooyr/log"
	"fmt"
	"io"
)

type music struct {
	voiceCons map[string]*discordgo.VoiceConnection
	sync.RWMutex
}

func (m *music) setVoiceCon(guildID string, voiceCon *discordgo.VoiceConnection) {
	m.Lock()
	m.voiceCons[guildID] = voiceCon
	m.Unlock()
}

func (m *music) getVoiceCon(guildID string) (*discordgo.VoiceConnection, bool) {
	m.RLock()
	voiceCon, ok := m.voiceCons[guildID]
	m.RUnlock()
	return voiceCon, ok
}

func newMusic() *music {
	return &music{voiceCons: make(map[string]*discordgo.VoiceConnection)}
}

func (m *music) join(ctx *command.Context) error {
	err := ctx.S.Scan()
	var guildID, channelID string
	switch err {
	case argument.ErrUnexpectedEnd:
		guildID, channelID, err = findAuthorVoiceChannel(ctx)
	case nil:
		name := ctx.S.Token()
		guildID, channelID, err = findVoiceChannel(ctx, name)
	default:
		return err
	}
	if err != nil {
		return err
	}

	// TODO do I need voice.ChangeChannel
	voiceCon, err := ctx.DG.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		return err
	}

	go func() {
		youtubeDL := exec.Command("youtube-dl", "-q", "-f", "bestaudio", "-o", "-", "https://youtu.be/dQw4w9WgXcQ")
		rickRoll, err := youtubeDL.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("volume=%f", 65/100)
		ffmpeg := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "quiet", "-i", "pipe:0",
			"-f", "data", "-map", "0:a", "-ar", "48k", "-ac", "2",
			"-af", fmt.Sprintf("volume=0.65"),
			"-acodec", "libopus", "-b:a", "128k", "pipe:1")
		ffmpeg.Stdin = rickRoll
		ffmpegOut, err := ffmpeg.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		err = youtubeDL.Start()
		if err != nil {
			log.Fatal(err)
		}
		defer log.Print("dooni")
		frames := make([][]byte, 0, 1000)
		var framesLock sync.Mutex
		go func() {
			for {
				if len(frames) > 0 {
					framesLock.Lock()
					frame := frames[0]
					frames = frames[1:]
					framesLock.Unlock()
					voiceCon.OpusSend <- frame
				}
			}
		}()
		err = ffmpeg.Start()
		if err != nil {
			log.Fatal(err)
		}
		// TODO do I need to do this?
		err = voiceCon.Speaking(true)
		if err != nil {
			log.Fatal(err)
		}
		for {
			p := make([]byte, 4000)
			n, err := ffmpegOut.Read(p)
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Fatal(err)
			}
			framesLock.Lock()
			frames = append(frames, p[:n])
			framesLock.Unlock()
		}
	}()
	m.setVoiceCon(guildID, voiceCon)
	return nil
}

// TODO https://github.com/bwmarrin/discordgo/wiki/FAQ#finding-a-users-voice-channel
func findAuthorGuild(ctx *command.Context) (*discordgo.Guild, error) {
	ch, err := ctx.DG.State.Channel(ctx.M.ChannelID)
	if err != nil {
		return nil, err
	}
	return ctx.DG.State.Guild(ch.GuildID)
}

func findAuthorVoiceChannel(ctx *command.Context) (string, string, error) {
	g, err := findAuthorGuild(ctx)
	if err != nil {
		return "", "", err
	}
	for _, vs := range g.VoiceStates {
		if vs.UserID == ctx.M.Author.ID {
			return g.ID, vs.ChannelID, nil
		}
	}
	return "", "", errors.New("Please join a voice channel or specify one.")
}

func findVoiceChannel(ctx *command.Context, name string) (string, string, error) {
	g, err := findAuthorGuild(ctx)
	if err != nil {
		return "", "", err
	}
	for _, ch := range g.Channels {
		if ch.Name == name {
			return g.ID, ch.ID, nil
		}
	}
	return "", "", errors.New("Requested voice channel not found.")
}

func (m *music) leave(ctx *command.Context) error {
	g, err := findAuthorGuild(ctx)
	if err != nil {
		return err
	}
	voiceCon, ok := m.getVoiceCon(g.ID)
	if !ok {
		return errors.New("Not connected.")
	}
	voiceCon.RLock()
	ready := voiceCon.Ready
	voiceCon.RUnlock()
	if !ready {
		return errors.New("Not connected.")
	}
	return voiceCon.Disconnect()
}
