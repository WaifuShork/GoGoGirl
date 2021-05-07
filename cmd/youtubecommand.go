package cmd

import (
	"GoGoGirl/framework"
	"bytes"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)


const resultFormat = "\n`%d` %s - %s (%s)"

var ytSessions YTSearchSessions = make(YTSearchSessions)

type (
	YTSearchSessions map[string]YTSearchSession

	YTSearchSession struct {
		Results []framework.YTSearchContent
	}
)

func ytSessionIdentifier(user *discordgo.User, channel *discordgo.Channel) string { 
	return user.ID + channel.ID
}

func formatDuration(input string) string { 
	return parseISO8601(input).String()
}

func YouTubeCommand(ctx framework.Context) { 
	if len(ctx.Args) == 0 { 
		ctx.Reply("Usage: `-youtube <search query>`")
		return
	}

	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID) 
	if sess == nil { 
		ctx.Reply("Not in a voice channel, to make the bot join one, use `-join`.")
		return
	}
	query := strings.Join(ctx.Args, " ")
	results, err := ctx.Youtube.Search(query)
	if err != nil { 
		ctx.Reply("An error occurred.")
		fmt.Println("Error searching youtube,", err)
		return
	}

	if len(results) == 0 { 
		ctx.Reply("No results found for your query `" + query + "`.")
		return
	}

	buffer := bytes.NewBufferString("__Search results__ for `" + query + "`:\n")
	for index, result := range results { 
		buffer.WriteString(fmt.Sprintf(resultFormat, index + 1, result.Title, result.ChannelTitle, formatDuration(result.Duration)))
	}

	buffer.WriteString("\n\nTo pick a song, use `-pick <number>`.")
	ytSessions[ytSessionIdentifier(ctx.User, ctx.TextChannel)] = YTSearchSession{results}
	ctx.Reply(buffer.String())
}