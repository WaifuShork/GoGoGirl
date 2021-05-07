package framework

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Represents a basic discord voice connection
// with Mutex channel lock
type Connection struct {
	VoiceConnection *discordgo.VoiceConnection
	Send 			chan []int16
	Lock            sync.Mutex // sync.Mutex
	SendPCM         bool
	StopRunning     bool
	Playing         bool
}

// Create a new connection
func NewConnection(voiceConnection *discordgo.VoiceConnection) *Connection { 
	connection := new(Connection)
	connection.VoiceConnection = voiceConnection
	return connection
}

// Disconnect 
func (connection Connection) Disconnect() { 
	connection.VoiceConnection.Disconnect()
}