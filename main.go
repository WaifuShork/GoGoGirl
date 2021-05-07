package main

import (
	"GoGoGirl/cmd"
	"GoGoGirl/framework"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	config     *framework.Config
	cmdHandler *framework.CommandHandler
	sessions   *framework.SessionManager
	youtube    *framework.YouTube
	botId      string
	botPrefix  string
)

func init() { 
	config = framework.LoadConfig("_config.json")
	botPrefix = config.Prefix
}

func main() {
	// setup framework shit
	cmdHandler = framework.NewCommandHandler()
	registerCommands()
	sessions = framework.NewSessionManager()
	youtube = &framework.YouTube{Config: config}
	
	// Create new discord bot instance
	discord, err := discordgo.New("Bot " + config.BotToken)
	if err != nil { 
		fmt.Println("Error creating discord session,", err)
		return
	}

	fmt.Println("Session created successfully")
	if config.UseSharding { 
		discord.ShardID = config.ShardId
		discord.ShardCount = config.ShardCount
	}

	// Get bot Id
	usr, err := discord.User("@me")
	if err != nil { 
		fmt.Println("Error obtaining account details,", err)
		return
	}

	botId = usr.ID; 
	// Add handlers for commands, and setup game status
	discord.AddHandler(commandHandler)
	//discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
	//	discord.UpdateGameStatus(0, config.DefaultStatus)
	//	guilds := discord.State.Guilds
	//	fmt.Println("Ready with", len(guilds), "guilds.")
	//})

	// Open a new discord websocket connection
	err = discord.Open()
	if err != nil { 
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Started")
	<-make(chan struct{})
}

// Responsible for handling incoming messages, and executing commands
func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author

	// Make sure the user isn't a bot for infinite loops
	if user.ID == botId || user.Bot { 
		return
	}

	// Make sure the prefix is correct
	content := message.Content
	if len(content) <= len(botPrefix) {
		return
	}

	// make sure the prefix matches
	if content[:len(botPrefix)] != botPrefix { 
		return
	}

	// slice the prefix off and return just the content
	content = content[len(botPrefix):]
	if len(content) < 1 {
		return
	}

	args := strings.Fields(content)
	name := strings.ToLower(args[0])
	command, found := cmdHandler.Get(name)
	if !found { 
		return
	}

	channel, err := discord.State.Channel(message.ChannelID) 
	if err != nil { 
		fmt.Println("Error getting channel,", err)
		return
	}

	guild, err := discord.State.Guild(channel.GuildID)
	if err != nil { 
		fmt.Println("Error getting guild,", err)
		return
	}

	ctx := framework.NewContext(discord, guild, channel, user, message, config, cmdHandler, sessions, youtube)
	ctx.Args = args[1:]
	c := *command
	c(*ctx)
}

// Handles checkers related reactions
func reactionsHandler(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	// Ignore all reactions created by the bot itself
	if reaction.UserID == session.State.SessionID  {
		return
	}

	// Retrieve extra information about the associated reaction
	msg, err := session.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	// Ignore reactions on messages that have an error or sent by bot
	if err != nil || msg == nil || msg.Author.ID != session.State.User.ID { 
		return
	}

	// Ignore messages that are not embeds with a command in the footer
	if len(msg.Embeds) != 1 || msg.Embeds[0].Footer == nil || msg.Embeds[0].Footer.Text == "" {
		return
	}

	


}

func registerCommands() { 
	cmdHandler.Register("help", cmd.HelpCommand, "Gives you this help message!")
	cmdHandler.Register("admin", cmd.AdminCommand, "Admin restricted command")
	cmdHandler.Register("join", cmd.JoinCommand, "Join a voice channel -join attic")
	cmdHandler.Register("leave", cmd.LeaveCommand, "Leaves current voice channel")
	cmdHandler.Register("play", cmd.PlayCommand, "Plays whats in the queue")
	cmdHandler.Register("stop", cmd.StopCommand, "Stops the music")
	cmdHandler.Register("add", cmd.AddCommand, "Add a song to the queue -add <youtube-link>")
	cmdHandler.Register("skip", cmd.SkipCommand, "Skip")
	cmdHandler.Register("queue", cmd.QueueCommand, "Print queue")
	cmdHandler.Register("eval", cmd.EvalCommand, "???")
	cmdHandler.Register("debug", cmd.DebugCommand, "???")
	cmdHandler.Register("clear", cmd.ClearCommand, "empty queue")
	cmdHandler.Register("current", cmd.CurrentCommand, "Name current song")
	cmdHandler.Register("youtube", cmd.YouTubeCommand, "???")
    cmdHandler.Register("shuffle", cmd.ShuffleCommand, "Shuffle queue")
    cmdHandler.Register("pausequeue", cmd.PauseCommand, "Pause song in place")
    cmdHandler.Register("pick", cmd.PickCommand, "Select a song from the list of YouTube results")
	cmdHandler.Register("say", cmd.SayCommand, "Say a message")
	cmdHandler.Register("encode", cmd.EncodeCommand, "Encodes a string to binary")
	cmdHandler.Register("decode", cmd.DecodeCommand, "decodes binary")
}