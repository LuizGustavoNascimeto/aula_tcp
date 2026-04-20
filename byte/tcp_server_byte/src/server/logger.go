package server

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	logFileName = "server_log.txt"
	logMutex    sync.Mutex
)

func AppLog(origin string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("[%s][%s] %s\n", timestamp, origin, msg)

	fmt.Print(line)

	logMutex.Lock()
	defer logMutex.Unlock()

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[%s][LOGGER] failed to open log file %s: %v\n", timestamp, logFileName, err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(line); err != nil {
		fmt.Printf("[%s][LOGGER] failed to write log file %s: %v\n", timestamp, logFileName, err)
	}
}
