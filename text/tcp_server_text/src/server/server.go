package server

import (
	"bufio"
	"crypto/sha512"
	"fmt"
	"net"
	"os"
	"strings"
	"tcp_server_text/src/session"
	"tcp_server_text/src/user"
)

type Server struct {
	addr     string
	users    *user.Repository
	sessions *session.Repository
}

const originServer = "servidor"

func NewServer(addr string) *Server {
	repo := user.NewRepository()
	hash := sha512.Sum512([]byte("password"))
	repo.Create(user.User{Login: "admin", Password: hash[:]})
	session := session.NewRepository()

	return &Server{
		addr:     addr,
		users:    repo,
		sessions: session,
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
	mode := os.Args[1]
	fmt.Println(mode)
	if mode == "byte" {
		AppLog(originServer, "AVERIGUANDO RESENHA")
	}

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			AppLog(originServer, "erro de leitura em %s: %v", conn.RemoteAddr().String(), err)
			return
		}

		s.handleMessage(conn, strings.TrimSpace(msg))

	}
}
