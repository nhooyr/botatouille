package music

import (
	"errors"

	"os/exec"
	"sync"

	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/botatouille/digo/argument"
	"github.com/nhooyr/botatouille/digo/command"
	"github.com/nhooyr/log"
	"runtime"
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
		ffmpeg := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "quiet", "-i", "pipe:0",
			"-f", "data", "-map", "0:a", "-ar", "48k", "-ac", "2",
			"-acodec", "libopus", "-b:a", "128k", "pipe:1")
		ffmpeg.Stdin = rickRoll
		ffmpegOut, err := ffmpeg.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		defer log.Print("dooni")
		framesChan := make(chan []byte, 100000)
		go func() {
			for {
				voiceCon.OpusSend <- <-framesChan
			}
		}()
		runtime.LockOSThread()
		err = youtubeDL.Start()
		if err != nil {
			log.Fatal(err)
		}
		err = ffmpeg.Start()
		if err != nil {
			log.Fatal(err)
		}
		for {
			// I read in the RFC that frames will not be bigger than this size
			p := make([]byte, 1275)
			n, err := ffmpegOut.Read(p)
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Fatal(err)
			}
			framesChan <- p[:n]
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
