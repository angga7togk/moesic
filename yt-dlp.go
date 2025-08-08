package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Youtube struct {
	Duration int64
	Url      string
}

func GetAudio(videoURL string) (*Youtube, error) {
	path := getYtDlpPath()
	cmd := exec.Command(path, "-f", "bestaudio[ext=m4a]", "-j", videoURL)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp error: %v", err)
	}

	var data struct {
		Duration int64    `json:"duration"`
		URL      string `json:"url"`
		Formats  []struct {
			Ext      string `json:"ext"`
			ACodec   string `json:"acodec"`
			VCodec   string `json:"vcodec"`
			FormatID string `json:"format_id"`
			URL      string `json:"url"`
		} `json:"formats"`
	}
	err = json.Unmarshal(output, &data)
	if err != nil {
		return nil, fmt.Errorf("json parse error: %v", err)
	}

	for _, format := range data.Formats {
		if format.Ext == "m4a" && format.VCodec == "none" {
			// * gw kurangin 1 soalnya biasanya, durasi 4:26 trus progrssnya cuman nyampe 4:25
			// * ya sementara solusinya cuman ini :)
			dur := data.Duration - 1 
			return &Youtube{
				Duration: dur,
				Url:      format.URL,
			}, nil
		}
	}

	return nil, fmt.Errorf("m4a audio format not found")
}
