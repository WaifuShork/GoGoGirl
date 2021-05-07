package cmd

import (
	"GoGoGirl/framework"
	"fmt"
)

func AddCommand(ctx framework.Context) { 
	if len(ctx.Args) == 0 { 
		ctx.Reply("Add command usage: `music add <song>`\nValid inputs: `youtube url`, `soundcloud url`, `youtube id`, `soundcloud id`")																	
		return
	}

	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess == nil { 
		ctx.Reply("Not in a voice channel, to make the bot join one, use `-join`")
		return
	}

	msg := ctx.Reply("Adding songs to queue...")
	for _, arg := range ctx.Args { 
		t, inp, err := ctx.Youtube.Get(arg)
	
		if err != nil {
			ctx.Reply("An error occurred!")
			fmt.Printf("error getting input, %s", err)
			return
		}

		switch t { 
		case framework.ErrorType:
			ctx.Reply("And error occurred!")
			fmt.Println("error type", t)
			return
		
		case framework.VideoType:
			video, err := ctx.Youtube.Video(*inp)
			if err != nil { 
				ctx.Reply("An error occurred!")
				fmt.Println("error getting video1,", err)
				return
			}
			song := framework.NewSong(video.Media, video.Title, arg)
			sess.Queue.Add(*song)
			ctx.Discord.ChannelMessageEdit(ctx.TextChannel.ID, msg.ID, "Added `" + song.Title + "` to the song queue." + 
			" Use `music play` to start playing the songs! To see the song queue, use `-queue`")
			// break
		case framework.PlaylistType:
			videos, err := ctx.Youtube.Playlist(*inp)
			if err != nil { 
				ctx.Reply("An error occurred!")
				fmt.Println("error getting playlist,", err)
				return
			}

			for _, v := range *videos { 
				id := v.Id
				_, i, err := ctx.Youtube.Get(id)
				if err != nil { 
					ctx.Reply("An error occurred!")
					fmt.Println("error getting video2,", err)
					continue
				}

				video, err := ctx.Youtube.Video(*i)
				if err != nil { 
					ctx.Reply("An error occurred!")
					fmt.Println("error getting video3,", err)
					return
				}

				song := framework.NewSong(video.Media, video.Title, arg)
				sess.Queue.Add(*song)
			}

			ctx.Reply("Finished adding songs to the playlist. Use `-play` to start playing the songs! " + 
			"To see the song queue, use `music queue`.")
			// break
		}
	}
}