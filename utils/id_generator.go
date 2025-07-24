package utils

import (
	"strconv"
	"strings"
	"time"
)

func IDGenerator() int {
	currentTime := time.Now()
	timestamp := currentTime.Format("2006-01-02 15:04:05.000")
	charsToRemove := [...]string{" ", "-", ":", "."}
	for _, char := range charsToRemove {
		timestamp = strings.ReplaceAll(timestamp, string(char), "")
	}
	intTs, _ := strconv.Atoi(timestamp)

	return intTs
}
