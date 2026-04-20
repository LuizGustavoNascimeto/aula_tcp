package server

import (
	"errors"
	"io"
	"net"
	"tcp_server_byte/src/datagram"
)

func (s *Server) WriteResponse(conn net.Conn, datagram *datagram.DatagramRes) {
	response := datagram.ToBytes()
	_, err := conn.Write(response)
	if err != nil {
		AppLog(conn.RemoteAddr().String(), "erro ao enviar resposta: %v", err)
	}
}
func (s *Server) logReadError(conn net.Conn, err error) {
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		AppLog(originServer, "cliente %s desconectou", conn.RemoteAddr().String())
		return
	}

	AppLog(originServer, "erro de leitura em %s: %v", conn.RemoteAddr().String(), err)
}
