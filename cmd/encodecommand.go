package cmd

import (
	"GoGoGirl/framework"
	"fmt"
)

func EncodeCommand(ctx framework.Context) { 
	if ctx.User.ID != ctx.Config.OwnerId { 
		return
	}

	if len(ctx.Args) == 0 {
		ctx.Reply("Usage: -encode <message>")
		return
	}

	var message string

	for _, contents := range ctx.Args { 
		message += fmt.Sprintf("%s ", contents)
	}
	
	strBin := binary(message)
	ctx.Reply(strBin)
}

func binary(s string) string {
    res := ""
    for _, c := range s {
        res = fmt.Sprintf("%s%.8b", res, c)
    }
	
    return res
}