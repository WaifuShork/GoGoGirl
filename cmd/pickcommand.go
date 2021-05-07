package cmd

import (
	"GoGoGirl/framework"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)


const invalidSongFormat = "Invalid song number `%d`. Min: 1, Max: %d"

func PickCommand(ctx framework.Context) { 
	argsLen := len(ctx.Args)
	if argsLen == 0 {
		ctx.Reply("Usage: `-pick <result number>`")
		return
	}
	if argsLen > 5 { 
		ctx.Reply("You cannot pick more than 5 songs at once.")
		return
	}
	identifier := ytSessionIdentifier(ctx.User, ctx.TextChannel)
	var ytSession YTSearchSession
	var ok bool
	if ytSession, ok = ytSessions[identifier]; !ok { 
		ctx.Reply("You haven't searched for a song yet")
		return
	}

	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID) 
	if sess == nil { 
		ctx.Reply("Not in a void channel, to make the bot join one, use `-join`.")
		return
	}

	rLen := len(ytSession.Results)
	var msg *discordgo.Message
	for i := 0; i < argsLen; i++ { 
		num, err := strconv.Atoi(ctx.Args[i])
		if err != nil { 
			ctx.Reply("An error occurred.")
			fmt.Print("Error parsing int,", err)
			return
		}

		if num < 1 || num > rLen { 
			ctx.Reply(fmt.Sprintf(invalidSongFormat, num, rLen))
			return
		}
		result := ytSession.Results[num - 1]

		_, inp, err := ctx.Youtube.Get(result.Id)
		video, err := ctx.Youtube.Video(*inp)
		song := framework.NewSong(video.Media, video.Title, result.Id)

		sess.Queue.Add(*song)
		if msg != nil { 
			msg, err = ctx.Discord.ChannelMessageEdit(ctx.TextChannel.ID, msg.ID, msg.Content + ", `" + song.Title + "`")
		} else { 
			msg = ctx.Reply("Added `" + song.Title + "`")
		}
	}

	if !sess.Queue.Running { 
		ctx.Discord.ChannelMessageEdit(ctx.TextChannel.ID, msg.ID, msg.Content + 
		" to the song queue.\nUse `-play` to start playing the sings. To see the song queue, use `-queue`.")
	}
}