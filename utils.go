package main

import (
	"fmt"
	"os"
	"strings"
)

func formatTime(seconds int64) string {
	min := seconds / 60
	sec := seconds % 60
	return fmt.Sprintf("%02d:%02d", min, sec)
}

func argsHas(ss ...string) bool {
	args := strings.Join(os.Args, ",")
	for _, s := range ss {
		if strings.Contains(args, s) {
			return true
		}
	}
	return false
}
