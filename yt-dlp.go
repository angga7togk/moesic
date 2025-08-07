package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Youtube struct {
	Duration float64
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
		Duration float64 `json:"duration"`
		URL      string  `json:"url"`
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
			dur := data.Duration
			if dur > 0.5 {
				dur -= 0.5
			}
			return &Youtube{
				Duration: dur,
				Url:      format.URL,
			}, nil
		}
	}

	return nil, fmt.Errorf("m4a audio format not found")
}
