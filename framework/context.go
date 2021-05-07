package framework

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Represents a Discord bot context for commands
type Context struct {
	Discord        *discordgo.Session
	Guild          *discordgo.Guild
	VoiceChannel   *discordgo.Channel
	TextChannel    *discordgo.Channel
	User           *discordgo.User
	Message        *discordgo.MessageCreate
	Args           []string

	Config         *Config
	CommandHandler *CommandHandler
	Sessions       *SessionManager
	Youtube        *YouTube
}

// Creates a new bot context
func NewContext(discord *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel, user *discordgo.User, 
				message *discordgo.MessageCreate, conf *Config, cmdHandler *CommandHandler, sessions *SessionManager, youtube *YouTube) (*Context) {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.Config = conf
	ctx.CommandHandler = cmdHandler
	ctx.Sessions = sessions
	ctx.Youtube = youtube
	return ctx
}

// Replies in the current channel that the command was executed
func (ctx Context) Reply(content string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.TextChannel.ID, content)
	if err != nil {
		fmt.Println("Error whilst sending message,", err)
		return nil
	}
	return msg
}

func (ctx *Context) GetVoiceChannel() *discordgo.Channel {
	if ctx.VoiceChannel != nil {
		return ctx.VoiceChannel
	}
	for _, state := range ctx.Guild.VoiceStates {
		if state.UserID == ctx.User.ID {
			channel, _ := ctx.Discord.State.Channel(state.ChannelID)
			ctx.VoiceChannel = channel
			return channel
		}
	}
	return nil
}