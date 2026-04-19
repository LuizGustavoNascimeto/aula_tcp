package commands

import (
	"tcp_server_text/src/session"
)

func EXIT(address string, sessionRepo *session.Repository) string {
	sessionRepo.Delete(address)
	return "BYE"
}
