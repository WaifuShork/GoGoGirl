package cmd

import "GoGoGirl/framework"

func LeaveCommand(ctx framework.Context) { 
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID) 
	if sess == nil { 
		ctx.Reply("Not in a void channel, to make the bot join one, use `-join`.")
		return
	}

	ctx.Sessions.Leave(ctx.Discord, *sess)
	ctx.Reply("Left <#" + sess.ChannelId + ">")
}