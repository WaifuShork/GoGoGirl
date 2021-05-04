package main

// Reminder:
// json.Marshal (encode): convert golang struct into json format.
// json.Unmarshal (decode): convert json into golang struct

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
}

var Token string

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {	
	run()
}


func run() { 
	config, err := readConfig("_config.json")
	if err != nil { 
		fmt.Println("error getting configuration.")
		return
	}

	discord, err := discordgo.New(config.Token)
	if err != nil {
		fmt.Println("error creating Discord sessions,", err)
		return
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID { 
		return
	}

	if message.Content == "ping" {
		fmt.Println(message.Content) 
		fmt.Println(message.ChannelID) 
		//return "pong"
		session.ChannelMessageSend(message.Content, "Pong")
	}

	if message.Content == "pong" { 
		fmt.Println(message.Content)
		fmt.Println(message.ChannelID) 
		//return "ping"
		session.ChannelMessageSend(message.ChannelID, "Pong")
	}
	//session.ChannelMessageSend(message.ChannelID, handleCommand(message.Content, "-")) 
}

func handleCommand(message, prefix string) (string) { 
	if strings.HasPrefix(message, prefix) {
		
	}
	return ""
}

func readConfig(cfgPath string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(cfgPath)

	if err != nil { 
		fmt.Println("File reading error", err)
		return config, err
	}

	// Deserialize json object
	err = json.Unmarshal(data, &config)
	if err != nil { 
		fmt.Println(err)
	}

	return config, err
}