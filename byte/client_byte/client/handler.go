package client

import (
	"net"

	"client_byte/commands"
	"client_byte/datagram"
)

func (c *Client) handleCommand(conn net.Conn, rawInput string) error {

	command, err := datagram.ParseReq(rawInput)
	if err != nil {
		return err
	}

	if err := c.sendRequest(conn, command); err != nil {
		return err
	}

	// Para ADDFILE, a resposta já foi lida em HandleAddFileReq
	// if command.CommandID == datagram.ADDFILE {
	// 	return nil
	// }

	resp, err := c.readResponse(conn)
	if err != nil {
		return err
	}

	// [5] Handle de resposta → aplica tratamento específico daquele comando
	return c.handleCommandResponse(conn, command, resp)
}

func (c *Client) handleCommandResponse(conn net.Conn, command *datagram.DatagramReq, resp *datagram.DatagramRes) error {
	var err error
	switch resp.CommandID {
	case datagram.GETFILESLIST:
		err = commands.HandleGetFilesListRes(conn, resp)
		if err != nil {
			return err
		}
	case datagram.GETFILE:
		err = commands.HandleGetFileRes(conn, command.Filename, resp)
		if err != nil {
			return err
		}
	}
	PrintResponse(resp)

	// Fallback para comandos sem handler específico
	ClientLog("comando %d finalizado", resp.CommandID)
	return nil
}
