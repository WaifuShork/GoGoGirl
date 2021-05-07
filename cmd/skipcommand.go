package cmd

import "GoGoGirl/framework"

func SkipCommand(ctx framework.Context) { 
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID) 
	if sess == nil { 
		ctx.Reply("Not in a voice channel, to make the bot join one, use `-join`.")
		return
	}

	sess.Stop()
	ctx.Reply("Skipped song.")
}