package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"tcp_server_byte/src/datagram"
	"time"
)

func AppLog(origin string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[%s][%s] %s\n", timestamp, origin, msg)
}

func (s *Server) WriteResponse(conn net.Conn, datagram *datagram.DatagramRes) {
	response := []byte{
		byte(datagram.MessageType),
		byte(datagram.CommandID),
		byte(datagram.StatusCode)}
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
