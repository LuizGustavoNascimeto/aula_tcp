package client

import (
	"fmt"
	"io"
	"net"

	"client_byte/commands"
	"client_byte/datagram"
)

// Envia arquivo após datagrama ADDFILE

// Lê resposta final após envio de arquivo

func (c *Client) sendRequest(conn net.Conn, command *datagram.DatagramReq) error {
	fileSize := uint32(0)
	if command.CommandID == datagram.ADDFILE {
		fileSize = command.GetFileSize()
	}

	payload, err := datagram.CreateReq(command.CommandID, command.Filename, fileSize)
	if err != nil {
		return err
	}

	if _, err = conn.Write(payload); err != nil {
		return err
	}

	if command.CommandID == datagram.ADDFILE {
		res1, err := c.readResponse(conn)
		if err != nil {
			return err
		}
		err = commands.HandleAddFileReq(conn, command.Filename, res1.StatusCode)
		if err != nil {
			return err
		}
		ClientLog("Arquivo %s transferido", command.Filename)
	}
	return nil
}

func (c *Client) readResponse(conn net.Conn) (*datagram.DatagramRes, error) {
	// [1] Leitura do comando (3 bytes de header)
	header := make([]byte, 3)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, fmt.Errorf("falha ao ler resposta do servidor: %w", err)
	}

	res, err := datagram.ParseRes(header)
	if err != nil {
		return nil, err
	}

	return res, nil
}
