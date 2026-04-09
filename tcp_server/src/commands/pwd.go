package commands

import (
	"tcp_server/src/session"
)

func PWD(address string, sessionRepo *session.Repository) string {
	userSession, exists := sessionRepo.Retrieve(address)
	if !exists {
		return "ERROR: NOT_AUTHENTICATED"
	}

	return userSession.CurrDir
}
