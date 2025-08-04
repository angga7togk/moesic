package main

import "fmt"

func formatTime(seconds float64) string {
	min := int(seconds) / 60
	sec := int(seconds) % 60
	return fmt.Sprintf("%02d:%02d", min, sec)
}