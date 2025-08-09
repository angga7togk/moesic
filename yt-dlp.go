package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)


func GetAudio(videoURL string) (string, error) {
	path := getYtDlpPath()
	cmd := exec.Command(path, "-f", "bestaudio[ext=m4a]", "-g", videoURL)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("yt-dlp error: %v", err)
	}
	url := strings.TrimSpace(out.String())
	return url, nil
}

