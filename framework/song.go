package framework

import (
	"os/exec"
	"strconv"
)

type Song struct {
	Media    string
	Title    string
	Duration *string
	Id       string
}

func (song Song) FFmpeg() *exec.Cmd { 
	return exec.Command("ffmpeg", "-i", song.Media, "s16le", "-ar", strconv.Itoa(FrameRate), "-ac", strconv.Itoa(Channels), "pipe:1")
}

func NewSong(media, title, id string) *Song { 
	song := new(Song)
	song.Media = media
	song.Title = title
	song.Id = id
	return song
}