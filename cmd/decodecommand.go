package cmd

import (
	"GoGoGirl/framework"
	"fmt"
	"strconv"
)

func DecodeCommand(ctx framework.Context) { 
	if ctx.User.ID != ctx.Config.OwnerId { 
		return
	}

	if len(ctx.Args) == 0 {
		ctx.Reply("Usage: -decode <message>")
		return
	}

	var message string

	for _, contents := range ctx.Args { 
		message += fmt.Sprintf("%s ", contents)
	}
	
	binary, err := strconv.ParseInt(message, 64, 8)

	if err != nil { 
		ctx.Reply(string(rune(binary)))
	}
}