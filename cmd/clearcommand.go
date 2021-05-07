package cmd

import "GoGoGirl/framework"

func ClearCommand(ctx framework.Context) { 
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess == nil { 
		ctx.Reply("Not in a voice channel, to make the boy join one, use `-join`.")
		return
	}

	if !sess.Queue.HasNext() {
		ctx.Reply("Queue is already empty.")
		return
	}

	sess.Queue.Clear()
	ctx.Reply("Cleared the song queue")
}