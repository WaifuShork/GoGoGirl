package framework

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	Channels  int = 2
	FrameRate int = 48000
	FrameSize int = 960
	MaxBytes  int = (FrameSize * 2) * 2
)

func (connection *Connection) sendPCM(voice *discordgo.VoiceConnection, pcm <-chan []int16) {
	connection.Lock.Lock()
	if connection.SendPCM || pcm == nil { 
		connection.Lock.Unlock()
		return
	}

	connection.SendPCM = true
	connection.Lock.Unlock()
	defer func() { 
		connection.SendPCM = false
	}()

	encoder, err := gopus.NewEncoder(FrameRate, Channels, gopus.Audio)
	if err != nil { 
		fmt.Println("NewEncoder error,", err)
		return
	}

	for { 
		receive, ok := <-pcm
		if !ok { 
			fmt.Println("PCM channel closed")
			return
		}

		opus, err := encoder.Encode(receive, FrameSize, MaxBytes)
		if err != nil { 
			fmt.Println("Encoding error,", err)
			return
		}

		if !voice.Ready || voice.OpusSend == nil {
			fmt.Printf("DiscordGo not ready for opus packets. %+v : %+v", voice.Ready, voice.OpusSend)
			return
		}

		voice.OpusSend <- opus
	}
}

func (connection *Connection) Play(ffmpeg *exec.Cmd) error { 
	if connection.Playing { 
		return errors.New("song already playing")
	}

	connection.StopRunning = false
	out, err := ffmpeg.StdoutPipe()
	if err != nil { 
		return err 
	}
	
	buffer := bufio.NewReaderSize(out, 16384)
	err = ffmpeg.Start()
	if err != nil { 
		return err
	}

	connection.Playing = true
	defer func() { 
		connection.Playing = false
	}()

	connection.VoiceConnection.Speaking(true)
	defer connection.VoiceConnection.Speaking(false)
	if connection.Send == nil { 
		connection.Send = make(chan []int16, 2)
	}

	go connection.sendPCM(connection.VoiceConnection, connection.Send)
	for { 
		if connection.StopRunning { 
			ffmpeg.Process.Kill()
			break
		}

		audioBuffer := make([]int16, FrameSize * Channels)
		err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF { 
			return nil
		}
		if err != nil {
			return err
		}
		connection.Send <- audioBuffer
	}
	return nil
}

func (connection *Connection) Stop() { 
	connection.StopRunning = true
	connection.Playing = false
}