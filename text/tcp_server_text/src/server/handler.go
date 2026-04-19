package server

import (
	"net"
	"strings"
	"tcp_server_text/src/commands"
)

func (s *Server) handleMessage(conn net.Conn, msg string) {
	AppLog(conn.RemoteAddr().String(), "%s", msg)
	command := strings.Split(msg, " ")

	if len(command) == 0 {
		s.WriteResponse(conn, "ERROR: INVALID_COMMAND")
		return
	}
	address := conn.RemoteAddr().String()
	var response string
	switch command[0] {
	case "CONNECT":
		response = commands.Connect(address, command, s.users, s.sessions)
	case "PWD":
		response = commands.PWD(address, s.sessions)
	case "CHDIR":
		response = commands.CHDIR(address, command, s.sessions)
	case "GETFILES":
		response = commands.GETFILES(address, s.sessions)
	case "GETDIRS":
		response = commands.GETDIRS(address, s.sessions)
	case "EXIT":
		response = commands.EXIT(address, s.sessions)
		s.WriteResponse(conn, response)
		conn.Close()
		return
	default:
		response = "ERROR: INVALID_COMMAND"
	}
	s.WriteResponse(conn, response)
}
