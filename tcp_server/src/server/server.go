package server

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net"
	"path"
	"path/filepath"
	"strings"
	"tcp_server/src/session"
	"tcp_server/src/user"
	"tcp_server/src/utils"
	"time"
)

type Server struct {
	addr    string
	repo    *user.Repository
	session *session.Repository
}

const originServer = "servidor"

func appLog(origin string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[%s][%s] %s\n", timestamp, origin, msg)
}

func (s *Server) writeResponse(conn net.Conn, payload string) {
	_, _ = conn.Write([]byte(payload + "\n"))
}

func NewServer(addr string) *Server {
	repo := user.NewRepository()
	hash := sha512.Sum512([]byte("password"))
	repo.Create(user.User{Login: "admin", Password: hash[:]})
	session := session.NewRepository()

	return &Server{
		addr:    addr,
		repo:    repo,
		session: session,
	}
}

func (s *Server) Run() error {
	appLog(originServer, "ouvindo em %s", s.addr)
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		appLog(originServer, "erro ao abrir porta %s: %v", s.addr, err)
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			appLog(originServer, "erro ao aceitar conexao: %v", err)
			return err
		}

		appLog(originServer, "nova conexao de %s", conn.RemoteAddr().String())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer appLog(originServer, "conexao encerrada com %s", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			appLog(originServer, "erro de leitura em %s: %v", conn.RemoteAddr().String(), err)
			return
		}

		s.handleMessage(conn, strings.TrimSpace(msg))
	}
}

func (s *Server) handleMessage(conn net.Conn, msg string) {
	appLog(conn.RemoteAddr().String(), "%s", msg)
	command := strings.Split(msg, " ")

	if len(command) == 0 {
		s.writeResponse(conn, "ERROR: INVALID_COMMAND")
	}

	if command[0] == "CONNECT" {
		s.commandConnect(conn, command)
		return
	}
	if command[0] == "PWD" {
		s.commandPWD(conn)
		return
	}
	if command[0] == "CHDIR" {
		s.commandCHDIR(conn, command)
		return
	}
}

func (s *Server) commandConnect(conn net.Conn, command []string) {
	if len(command) < 3 {
		appLog(originServer, "CONNECT invalido de %s", conn.RemoteAddr().String())
		s.writeResponse(conn, "ERROR: INVALID_COMMAND")
		return
	}

	pass, err := hex.DecodeString(command[2])
	if err != nil {
		appLog(originServer, "hash invalido recebido de %s para usuario %s", conn.RemoteAddr().String(), command[1])
		s.writeResponse(conn, "ERROR: INVALID_HASH")
		return
	}

	if s.repo.HandleAuth(command[1], pass) {
		rootDir := fmt.Sprintf("/%s", command[1])
		s.session.Create(session.Session{ID: conn.RemoteAddr().String(), CurrDir: rootDir, User: command[1], RootDir: rootDir})
		appLog(originServer, "usuario %s autenticado com sucesso (%s)", command[1], conn.RemoteAddr().String())
		s.writeResponse(conn, "SUCCESS")
		return
	}

	appLog(originServer, "falha de autenticacao para usuario %s (%s)", command[1], conn.RemoteAddr().String())

	s.writeResponse(conn, "ERROR")
}

func (s *Server) commandPWD(conn net.Conn) {
	session, exists := s.session.Retrieve(conn.RemoteAddr().String())
	if !exists {
		appLog(originServer, "PWD negado para %s: usuario nao autenticado", conn.RemoteAddr().String())
		s.writeResponse(conn, "ERROR: NOT_AUTHENTICATED")
		return
	}
	appLog(originServer, "PWD de %s (%s): %s", session.User, conn.RemoteAddr().String(), session.CurrDir)
	s.writeResponse(conn, session.CurrDir)
}

// recebe um caminho "global"
func (s *Server) commandCHDIR(conn net.Conn, command []string) {
	if len(command) < 2 {
		appLog(originServer, "CHDIR invalido de %s", conn.RemoteAddr().String())
		s.writeResponse(conn, "ERROR")
		return
	}
	userSession, exists := s.session.Retrieve(conn.RemoteAddr().String())
	if !exists {
		appLog(originServer, "CHDIR negado para %s: usuario nao autenticado", conn.RemoteAddr().String())
		s.writeResponse(conn, "ERROR")
		return
	}
	newDir := command[1]
	if !strings.HasPrefix(newDir, "/") {
		newDir = "/" + newDir
	}
	if !isValidDir(userSession.RootDir, newDir) {
		appLog(originServer, "CHDIR invalido para usuario %s: %s", userSession.User, newDir)
		s.writeResponse(conn, "ERROR: INVALID_DIRECTORY")
		return
	}
	s.session.Update(userSession.ID, newDir)
	appLog(originServer, "usuario %s alterou diretorio para %s", userSession.User, newDir)
	s.writeResponse(conn, "SUCCESS")

}

func isValidDir(root string, newDir string) bool {
	// root/newDir sao caminhos logicos (com '/'), por isso usamos path.Clean
	cleanRoot := path.Clean("/" + strings.TrimPrefix(root, "/"))
	cleanNewDir := path.Clean("/" + strings.TrimPrefix(newDir, "/"))

	if cleanNewDir != cleanRoot && !strings.HasPrefix(cleanNewDir, cleanRoot+"/") {
		return false
	}

	relativeDir := filepath.FromSlash(strings.TrimPrefix(cleanNewDir, "/"))
	targetDir := filepath.Join("users_files", relativeDir)

	return utils.DirExists(targetDir)
}
