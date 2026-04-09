package commands

import (
	"tcp_server/src/session"
	"tcp_server/src/utils"
)

func CHDIR(address string, command []string, sessionRepo *session.Repository) string {
	if len(command) < 2 {
		return "ERROR: INVALID_COMMAND"
	}
	userSession, exists := sessionRepo.Retrieve(address)
	if !exists {
		return "ERROR: NOT_AUTHENTICATED"
	}
	newDir := command[1]
	if !utils.IsValidDir(userSession.RootDir, newDir) {
		return "ERROR: INVALID_DIRECTORY"
	}
	sessionRepo.Update(userSession.ID, newDir)
	return "SUCCESS"
}
