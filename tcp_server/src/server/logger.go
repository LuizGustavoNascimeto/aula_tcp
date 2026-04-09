package server

import (
	"fmt"
	"net"
	"time"
)

func AppLog(origin string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[%s][%s] %s\n", timestamp, origin, msg)
}

func (s *Server) WriteResponse(conn net.Conn, payload string) {
	_, _ = conn.Write([]byte(payload + "\n"))
}
