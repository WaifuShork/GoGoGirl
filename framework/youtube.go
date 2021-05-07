package framework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

type VideoResponse struct {
	Formats []struct {
		Url string `json:"url"`
	} `json:"formats"`
	Title string `json:"title"`
}

type VideoResult struct {
	Media string
	Title string
}

type PlaylistVideo struct {
	Id string `json:"id"`
}

type YTSearchContent struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"descrption"`
	ChannelTitle string `json:"channel_title"`
	Duration     string `json:"duration"`
}

type YTAPIResponse struct {
	Error   bool              `json:"error"`
	Content []YTSearchContent `json:"content"`
}

type YouTube struct {
	Config *Config
}

const (
	ErrorType    = -1
	VideoType    = 0
	PlaylistType = 1
)

// Get the type of the video/url that's been passed
func (youtube YouTube) getType(input string) int {
	if strings.Contains(input, "upload_date") {
		return VideoType
	}
	if strings.Contains(input, "_type") {
		return PlaylistType
	}
	return ErrorType
}

func (youtube YouTube) Get(input string) (int, *string, error) { 
	cmd := exec.Command("youtube-dl", "--skip-download", "--print-json", "--flat-playlist", input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ErrorType, nil, err
	}
	str := out.String()
	return youtube.getType(str), &str, nil
}

func (youtube YouTube) Video(input string) (*VideoResult, error) { 
	var response VideoResponse
	err := json.Unmarshal([]byte(input), &response)
	if err != nil { 
		return nil, err
	}

	return &VideoResult{response.Formats[0].Url, response.Title}, nil
}

func (youtube YouTube) Playlist(input string) (*[]PlaylistVideo, error) { 
	lines := strings.Split(input, "\n")
	videos := make([]PlaylistVideo, 0)
	for _, line := range lines { 
		if len(line) == 0 {
			continue
		}

		var video PlaylistVideo
		fmt.Println("line,", line)
		err := json.Unmarshal([]byte(line), &video)
		if err != nil { 
			return nil, err
		}
		videos = append(videos, video)
	}

	return &videos, nil
}

func (youtube YouTube) buildUrl(query string) (*string, error) { 
	base := youtube.Config.ServiceUrl + "/v1/youtube/search"
	address, err := url.Parse(base)
	if err != nil { 
		return nil, err
	}

	params := url.Values{}
	params.Add("search", query)
	address.RawQuery = params.Encode()
	str := address.String()
	return &str, nil
}

func (youtube YouTube) Search(query string) ([]YTSearchContent, error) { 
	address, err := youtube.buildUrl(query)
	if err != nil { 
		return nil, err
	}

	response, err := http.Get(*address)
	if err != nil { 
		return nil, err
	}

	var apiResponse YTAPIResponse
	json.NewDecoder(response.Body).Decode(&apiResponse)
	return apiResponse.Content, nil
}