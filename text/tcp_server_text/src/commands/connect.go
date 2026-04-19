package commands

import (
	"encoding/hex"
	"fmt"
	"strings"
	"tcp_server_text/src/session"
	"tcp_server_text/src/user"
)

func Connect(address string, command []string, userRepo *user.Repository, sessionRepo *session.Repository) string {
	if len(command) < 2 {
		return "ERROR: INVALID_COMMAND"
	}
	params := strings.Split(command[1], ",")
	pass, err := hex.DecodeString(params[1])
	if err != nil {

		return "ERROR: INVALID_HASH"
	}
	if userRepo.HandleAuth(params[0], pass) {
		rootDir := fmt.Sprintf("/%s", params[0])
		sessionRepo.Create(session.Session{ID: address, CurrDir: rootDir, User: params[0], RootDir: rootDir})
		return "SUCCESS"
	}
	return "ERROR"
}
