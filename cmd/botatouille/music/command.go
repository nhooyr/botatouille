package music

import (
	"github.com/nhooyr/botatouille/digo/command"
)

var Command = command.NewRouter("m", "music")

func init() {
	m := newMusic()
	join := command.NewCommand("join", "[channel]", "join", m.join)
	leave := command.NewCommand("leave", "", "leave", m.leave)

	Command.Append(join)
	Command.Append(leave)
}