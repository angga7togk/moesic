package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func getDuration(url string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries",
		"format=duration", "-of", "default=noprint_wrappers=1:nokey=1", url)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	str := strings.TrimSpace(string(out))
	return strconv.ParseFloat(str, 64)
}

func play(url string) *exec.Cmd {
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", url)
	err := cmd.Start()
	if err != nil {
		fmt.Println("Play error:", err)
	}
	return cmd
}
