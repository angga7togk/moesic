package main

import "fmt"

func formatTime(seconds int64) string {
	min := seconds / 60
	sec := seconds % 60
	return fmt.Sprintf("%02d:%02d", min, sec)
}
