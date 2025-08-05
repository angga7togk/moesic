package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func play(url string) *exec.Cmd {
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "info", "-infbuf", url)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("stderr error:", err)
		return nil
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Play error:", err)
		return nil
	}

	go func() {
		reader := bufio.NewReader(stderrPipe)
		for {
			line, err := reader.ReadString('\r')
			if err != nil {
				break
			}

			if strings.Contains(line, "A-V:") || strings.Contains(line, "M-A:") || strings.Contains(line, "M-V:") {
				fields := strings.Fields(line)
				if len(fields) > 0 {
					timestampStr := fields[0]
					timestamp, err := strconv.ParseFloat(timestampStr, 64)
					if err == nil {
						currentTime = timestamp
					}
				}
			}
		}
	}()

	return cmd
}
