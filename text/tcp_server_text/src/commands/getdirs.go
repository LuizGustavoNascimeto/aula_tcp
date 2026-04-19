package commands

import (
	"fmt"
	"strings"
	"tcp_server_text/src/session"
	"tcp_server_text/src/utils"
)

func GETDIRS(address string, sessionRepo *session.Repository) string {
	userSession, exists := sessionRepo.Retrieve(address)
	if !exists {
		return "ERROR: NOT_AUTHENTICATED"
	}
	dirs, err := utils.ListDirs(userSession.CurrDir)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	if len(dirs) == 0 {
		return "NOT_FOUND"
	}
	resp := strings.Join(dirs, "\n")
	size := fmt.Sprint(len(dirs)) + "\n"
	return size + resp
}
