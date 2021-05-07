package framework

import "github.com/bwmarrin/discordgo"

type Session struct {
	Queue      *SongQueue
	GuildId    string
	ChannelId  string
	Connection *Connection
}

type SessionManager struct {
	Sessions map[string]*Session
}

type JoinProperties struct {
	Muted    bool
	Deafened bool
}

func NewSession(guildId, channelId string, connection *Connection) *Session {
	session := new(Session)
	session.Queue = newSongQueue()
	session.GuildId = guildId
	session.ChannelId = channelId
	session.Connection = connection
	return session
}

func (session Session) Play(song Song) error {
	return session.Connection.Play(song.FFmpeg())
}

func (session *Session) Stop() {
	session.Connection.Stop()
}

func NewSessionManager() *SessionManager {
	return &SessionManager{make(map[string]*Session)}
}

func (manager SessionManager) GetByGuild(guildId string) *Session {
	for _, sess := range manager.Sessions {
		if sess.GuildId == guildId {
			return sess
		}
	}
	return nil
}

func (manager SessionManager) GetByChannel(channelId string) (*Session, bool) {
	sess, found := manager.Sessions[channelId]
	return sess, found
}

func (manager *SessionManager) Join(discord *discordgo.Session, guildId, channelId string, properties JoinProperties) (*Session, error) { 
	vc, err := discord.ChannelVoiceJoin(guildId, channelId, properties.Muted, properties.Deafened)
	if err != nil { 
		return nil, err
	}

	sess := NewSession(guildId, channelId, NewConnection(vc))
	manager.Sessions[channelId] = sess
	return sess, nil
}

func (manager *SessionManager) Leave(discord *discordgo.Session, session Session) { 
	session.Connection.Stop()
	session.Connection.Disconnect()
	delete(manager.Sessions, session.ChannelId)
}