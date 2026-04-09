package commands

import (
	"fmt"
	"strings"
	"tcp_server/src/session"
	"tcp_server/src/utils"
)

func GETFILES(address string, sessionRepo *session.Repository) string {
	userSession, exists := sessionRepo.Retrieve(address)
	if !exists {
		return "ERROR: NOT_AUTHENTICATED"
	}
	files, err := utils.ListFiles(userSession.CurrDir)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	if len(files) == 0 {
		return "NOT_FOUND"
	}
	resp := strings.Join(files, "\n")
	size := fmt.Sprint(len(files)) + "\n"
	return size + resp

}
