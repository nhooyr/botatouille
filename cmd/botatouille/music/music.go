package music

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/botatouille/digo/argument"
	"github.com/nhooyr/botatouille/digo/command"
	"errors"
)

type music struct {
	guildMusics map[string]*guildMusic
	sync.RWMutex
}

func newMusic() *music {
	return &music{guildMusics: make(map[string]*guildMusic)}
}

func (m *music) setGuildMusic(guildID string, gm *guildMusic) {
	m.Lock()
	m.guildMusics[guildID] = gm
	m.Unlock()
}

func (m *music) getGuildMusic(guildID string) (*guildMusic, bool) {
	m.RLock()
	gm, ok := m.guildMusics[guildID]
	m.RUnlock()
	return gm, ok
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

	voiceCon, err := ctx.DG.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		return err
	}

	gm := newGuildMusic(voiceCon)
	m.setGuildMusic(guildID, gm)
	return nil
}

func findAuthorGuild(ctx *command.Context) (*discordgo.Guild, error) {
	for _, g := range ctx.DG.State.Guilds {
		for _, c := range g.Channels {
			if c.ID == ctx.M.ChannelID {
				return g, nil
			}
		}
	}
	// TODO crazy
	return nil, errors.New("Unable to find your guild. Crazy, please report to Anmol.")
}

func findAuthorVoiceChannel(ctx *command.Context) (string, string, error) {
	// One nuance of this method is that it allows the bot to join another guild's
	// voice connection when being controlled from a different guild if you are there.
	for _, g := range ctx.DG.State.Guilds {
		for _, vs := range g.VoiceStates {
			if vs.UserID == ctx.M.Author.ID {
				return g.ID, vs.ChannelID, nil
			}
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
		if ch.Name == name && ch.Type == "voice" {
			return g.ID, ch.ID, nil
		}
	}
	return "", "", errors.New("Requested voice channel not found.")
}

func (m *music) findGuildMusic(ctx *command.Context) (*guildMusic, error) {
	g, err := findAuthorGuild(ctx)
	if err != nil {
		return nil, err
	}
	gm, ok := m.getGuildMusic(g.ID)
	if !ok {
		return nil, errors.New("Not connected.")
	}
	return gm, nil
}

func (m *music) leave(ctx *command.Context) error {
	gm, err := m.findGuildMusic(ctx)
	if err != nil {
		return err
	}
	return gm.stop()
}

func (m *music) add(ctx *command.Context) error {
	gm, err := m.findGuildMusic(ctx)
	if err != nil {
		return err
	}
	gm.q.append(&video{id: "dQw4w9WgXcQ"})
	return nil
}
