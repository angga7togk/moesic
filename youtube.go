package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

type Youtube struct {
	Duration float64
	Url      string
}

func GetYoutube(videoURL string) (*Youtube, error) {
	var binary string
	switch runtime.GOOS {
	case "linux":
		binary = "./bin/yt-dlp_linux"
	case "windows":
		binary = "./bin/yt-dlp_win.exe"
	case "darwin":
		binary = "./bin/yt-dlp_macos"
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	cmd := exec.Command(binary, "-f", "bestaudio[ext=m4a]", "-j", videoURL)
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

			// kurangin 0.5 biar pas progress nya ke 100% dan menghindari macet saat next music
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
