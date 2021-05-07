package cmd

import (
	"GoGoGirl/framework"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)


type JsonStructure struct {
	Entries []string
}

const todoFile = "todo.json"

func (s *JsonStructure) add(entry string) {
	s.Entries = append(s.Entries, entry)
}

func readTodo() *JsonStructure {
	bBody, err := ioutil.ReadFile(todoFile)

	if err != nil { 
		fmt.Println("error reading todo file,", err)
		return nil
	}

	s := &JsonStructure{make([]string, 0)}
	err = json.Unmarshal(bBody, s)
	if err != nil { 
		fmt.Println("error unmarshalling todo file,", err)
		return nil
	}
	return s
}

func writeTodo(s JsonStructure) { 
	body, err := json.Marshal(s)
	if err != nil { 
		fmt.Println("error marshalling todo file,", err)
		return
	}

	err = ioutil.WriteFile(todoFile, body, os.ModeAppend)
	if err != nil { 
		fmt.Println("error writing todo file,", err)
		return
	}
}

func AdminCommand(ctx framework.Context) { 
	if ctx.User.ID != ctx.Config.OwnerId { 
		return
	}

	if len(ctx.Args) == 0 {
		ctx.Reply("Usage: -admin <subcommand>\nSubcommands: stop")
		return
	}

	switch strings.ToLower(ctx.Args[0]) { 
	case "stop":
		ctx.Reply("Bye :wave:")
		err := ctx.Discord.Close()
		if err != nil {
			return 
		}
		os.Exit(-1)
	case "todo":
		{
			str := readTodo()
			buffer := bytes.NewBufferString("todo list entries:")
			for i, s := range str.Entries { 
				buffer.WriteString(fmt.Sprintf("\n`%02d` %s", i + 1, s))
			}	
			ctx.Reply(buffer.String())
			break
		}
	case "addTodo":
		{
			entry := strings.Join(ctx.Args[1:], " ")
			str := readTodo()
			if str == nil { 
				str = &JsonStructure{make([]string, 0)}
			}
			str.add(entry)
			writeTodo(*str)
			ctx.Reply("wrote todo list")
			break
		}
	default:
		ctx.Reply("Invalid subcommand!")
	}
}