package cmd

import (
	"GoGoGirl/framework"
	"bytes"
)


func DebugCommand(ctx framework.Context) { 
	if ctx.Config.OwnerId != ctx.User.ID { 
		return
	}

	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess == nil { 
		ctx.Reply("no current session")
		return
	}

	queue := sess.Queue
	q := queue.Get()
	buffer := bytes.Buffer{}
	for _, song := range q { 
		buffer.WriteString(song.Id + " ")
	}
	ctx.Reply(buffer.String())
}