package cmd

import "GoGoGirl/framework"

func PauseCommand(ctx framework.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID) 

	if sess == nil { 
		ctx.Reply("Not in a voice channel, to make the bot join one, use `-join`.")
		return
	}

	queue := sess.Queue
	if !queue.HasNext() { 
		ctx.Reply("Queue is empty. Add songs with `-add`.")
		return
	}

	queue.Pause()
	ctx.Reply("The queue has paused and will stop playing after this song. To resume the queue, use `-play`.")
}