package commands

import (
	"client_byte/datagram"
	"client_byte/utils"
	"fmt"
	"net"
)

func HandleAddFileReq(conn net.Conn, filename string, status int8) error {
	if status != datagram.STATUS_SUCCESS {
		return fmt.Errorf("ADDFILE recusado no primeiro status")
	}
	// Lê resposta final após envio do arquivo
	if err := utils.SendFile(conn, filename); err != nil {
		return err
	}
	return nil
}
