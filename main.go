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

	AdministratorId uint64 `json:"administratorId"` 
	ModeratorId     uint64 `json:"moderatorId"`
}

type Command struct { 
	Name  string
	Value interface{}
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

	config, err := readConfig("_config.json")
	if err != nil { 
		fmt.Println("unable to parse config for command context.")
		return
	}

	if message.Author.ID == session.State.User.ID { 
		return
	}

	msg := parseCommand(message.Content, config)
	session.ChannelMessageSend(message.ChannelID, msg) 
}

func parseCommand(message string, config Config) (string) { 
	commands := make(map[string]Command)

	commands["ping"] = Command{Name: "ping", Value: "pong"}
	commands["pong"] = Command{Name: "pong", Value: "ping"}
	commands["modId"] = Command{Name: "modId", Value: config.ModeratorId}
	commands["adId"] = Command{Name: "adId", Value: config.AdministratorId}

	// Trim the prefix since we only need the raw message contents.
	noPrefix := strings.TrimPrefix(message, config.Prefix)

	for range commands { 
		if commands[noPrefix].Name == noPrefix { 
			fmt.Println(commands[noPrefix].Value)
			return fmt.Sprintf("%v", commands[noPrefix].Value)
		}
	}

	// sprintf the value since it's 'interface{}'.
	value := fmt.Sprintf("%v", commands[noPrefix].Value)
	fmt.Println("Command Name: " + commands[noPrefix].Name, "Command Value: " + value)

	// return the users message in full without the prefix, for error reporting. 
	return noPrefix;
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