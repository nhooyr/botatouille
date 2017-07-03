package main

import "github.com/nhooyr/botatouille/command"

func custom_command(ctx *command.Context) {
	ctx.DG.ChannelVoiceJoin()
}
