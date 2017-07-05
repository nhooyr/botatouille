package music

import (
	"errors"

	"log"
	"os/exec"
	"sync"

	"bufio"
	"encoding/binary"
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/botatouille/digo/argument"
	"github.com/nhooyr/botatouille/digo/command"
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

		err = youtubeDL.Start()
		if err != nil {
			log.Fatal(err)
		}

		// Create opus stream
		stream, err := convertToOpus(rickRoll)
		if err != nil {
			log.Fatal(err)
		}

		// Loop until the audio is done playing.
		for {
			opus, err := readOpus(stream)
			if err != nil {
				if err == io.ErrUnexpectedEOF || err == io.EOF {
					return
				}
				log.Print(err)
				return
			}

			voiceCon.OpusSend <- opus
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

// Reads an opus packet to send over the vc.OpusSend channel
func readOpus(source io.Reader) ([]byte, error) {
	var opuslen int16
	err := binary.Read(source, binary.LittleEndian, &opuslen)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, err
		}
		return nil, errors.New("ERR reading opus header")
	}

	var opusframe = make([]byte, opuslen)
	err = binary.Read(source, binary.LittleEndian, &opusframe)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, err
		}
		return nil, errors.New("ERR reading opus frame")
	}

	return opusframe, nil
}

// convertToOpus converts the given io.Reader stream to an Opus stream
// Using ffmpeg and dca-rs
func convertToOpus(rd io.Reader) (io.Reader, error) {

	// Convert to a format that can be passed to dca-rs
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpeg.Stdin = rd
	ffmpegout, err := ffmpeg.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// Convert to opus
	// TODO maybe use dca-rs later? https://github.com/bwmarrin/dca
	dca := exec.Command("dca", "--raw", "-i", "pipe:0")
	dca.Stdin = ffmpegout
	dcaout, err := dca.StdoutPipe()
	dcabuf := bufio.NewReaderSize(dcaout, 1024)
	if err != nil {
		return nil, err
	}

	// Start ffmpeg
	err = ffmpeg.Start()
	if err != nil {
		return nil, err
	}

	// Start dca-rs
	err = dca.Start()
	if err != nil {
		return nil, err
	}

	// Returns a stream of opus data
	return dcabuf, nil
}