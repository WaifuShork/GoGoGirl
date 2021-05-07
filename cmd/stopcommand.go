package cmd

import "GoGoGirl/framework"

func StopCommand(ctx framework.Context) { 
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID) 

	if sess == nil {
		ctx.Reply("Not in a voice channel, to make the bot join one, use `-join`.")
		return
	}

	if sess.Queue.HasNext() {
		sess.Queue.Clear()
	}
	sess.Stop()
}