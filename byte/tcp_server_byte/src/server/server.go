package server

import (
	"bufio"
	"io"
	"net"
	"tcp_server_byte/src/buffer"
	"tcp_server_byte/src/datagram"
)

type Server struct {
	addr string
}

const originServer = "servidor"

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Run() error {
	AppLog(originServer, "ouvindo em %s", s.addr)
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		AppLog(originServer, "erro ao abrir porta %s: %v", s.addr, err)
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			AppLog(originServer, "erro ao aceitar conexao: %v", err)
			return err
		}

		AppLog(originServer, "nova conexao de %s", conn.RemoteAddr().String())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer AppLog(originServer, "conexao encerrada com %s", conn.RemoteAddr().String())
	defer buffer.GetReaderPool().Remove(conn)

	reader := buffer.GetReaderPool().Reader(conn)
	for {
		message, err := s.readRequest(reader)
		if err != nil {
			s.logReadError(conn, err)
			return
		}
		s.handleMessage(conn, reader, message)
	}
}

func (s *Server) readRequest(reader *bufio.Reader) (*datagram.DatagramReq, error) {
	// Header fixo: tipo, comando, tamanho do filename.
	header := make([]byte, 3)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, err
	}

	message, err := datagram.ParseReq(header)
	if err != nil {
		return nil, err
	}

	message.HandleFilename(reader)
	return message, nil
}
