package music

import (
	"github.com/nhooyr/botatouille/digo/command"
)

// TODO turn into function, no init pls
var Command = command.NewRouter("m", "music")

func init() {
	m := newMusic()
	join := command.NewCommand("join", "[channel]", "join", m.join)
	leave := command.NewCommand("leave", "", "leave", m.leave)
	add := command.NewCommand("add", "", "add", m.add)

	Command.Append(join)
	Command.Append(leave)
	Command.Append(add)
}