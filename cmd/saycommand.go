package cmd

import (
	"GoGoGirl/framework"
	"fmt"
)

func SayCommand(ctx framework.Context) { 
	if ctx.User.ID != ctx.Config.OwnerId { 
		return
	}

	if len(ctx.Args) == 0 {
		ctx.Reply("Usage: -say <message>")
		return
	}

	var message string

	for _, contents := range ctx.Args { 
		message += fmt.Sprintf("%s ", contents)
	}

	ctx.Reply(message)
}