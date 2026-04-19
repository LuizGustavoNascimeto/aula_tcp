package server

import (
	"bufio"
	"net"
	"tcp_server_byte/src/commands"
	"tcp_server_byte/src/datagram"
)

func (s *Server) handleMessage(conn net.Conn, reader *bufio.Reader, msg *datagram.DatagramReq) {
	AppLog(conn.RemoteAddr().String(), "%s", msg.String())

	//address := conn.RemoteAddr().String()
	var res datagram.DatagramRes
	switch msg.CommandID {
	case datagram.ADDFILE:
		//esse resposta infica que o servidor está pronto para receber o arquivo
		s.WriteResponse(conn, datagram.CreateRes(datagram.ADDFILE, datagram.STATUS_SUCCESS))
		status, err := commands.ADDFILE(reader, msg.Filename, msg.GetFileSize())
		if err != nil {
			AppLog(conn.RemoteAddr().String(), "erro ao processar comando ADDFILE: %v", err)
		}
		res = *datagram.CreateRes(datagram.ADDFILE, status)
	// case "PWD":
	// 	response = commands.PWD(address, s.sessions)
	// case "CHDIR":
	// 	response = commands.CHDIR(address, command, s.sessions)
	// case "GETFILES":
	// 	response = commands.GETFILES(address, s.sessions)
	// case "GETDIRS":
	// 	response = commands.GETDIRS(address, s.sessions)
	// case "EXIT":
	// 	response = commands.EXIT(address, s.sessions)
	// 	s.WriteResponse(conn, response)
	// 	conn.Close()
	// 	return
	default:
		res = *datagram.CreateRes(datagram.INVALID_COMMAND, datagram.STATUS_ERROR)
	}
	s.WriteResponse(conn, &res)
}
