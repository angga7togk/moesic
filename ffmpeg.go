package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func durationToSeconds(s string) (int64, error) {
	var h, m int
	var sec float64
	_, err := fmt.Sscanf(s, "%02d:%02d:%f", &h, &m, &sec)
	if err != nil {
		return 0, err
	}
	total := int64(h)*3600 + int64(m)*60 + int64(sec)
	return total, nil
}

var (
	globalCurrentTime     int64 = 0 // * global progres music player
	globalCurrentDuration int64 = 0
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

			if strings.Contains(line, "Duration:") {
				// * ambil string setelah kata "Duration:"
				parts := strings.SplitN(line, "Duration:", 2)
				if len(parts) > 1 {
					after := strings.TrimSpace(parts[1]) // ? "00:03:23.96, start: 0.000000, bitrate: 129 kb/s"
					// * potong sampai koma pertama
					durationStr := strings.SplitN(after, ",", 2)[0]
					f, err := durationToSeconds(durationStr)
					if err == nil {
						globalCurrentDuration = f
					}
				}
			}

			// * buat dengerin progress audio nya, entah itu video audio apapu itu 
			if strings.Contains(line, "A-V:") || strings.Contains(line, "M-A:") || strings.Contains(line, "M-V:") {
				fields := strings.Fields(line)
				if len(fields) > 0 {
					timestampStr := fields[0]
					f, err := strconv.ParseFloat(timestampStr, 64)
					if err == nil {
						globalCurrentTime = int64(f)
					}
				}
			}
		}
	}()

	return cmd
}
