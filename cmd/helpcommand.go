package cmd

import (
	"GoGoGirl/framework"
	"bytes"
	"fmt"
) 

func HelpCommand(ctx framework.Context) {
	cmds := ctx.CommandHandler.GetCommands()
	buffer := bytes.NewBufferString("Commands: \n")
	for cmdName, cmdStruct := range cmds { 
		if len(cmdName) == 1 { 
			continue
		}
		msg := fmt.Sprintf("\t %s%s - %s\n", ctx.Config.Prefix, cmdName, cmdStruct.GetHelp())
		buffer.WriteString(msg)
	}

	str := buffer.String()
	ctx.Reply(str[:len(str) - 2])
}