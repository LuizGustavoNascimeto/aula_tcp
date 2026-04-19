package commands

import (
	"fmt"
	"net"

	"client_byte/datagram"
)

func HandleADDFILE(
	conn net.Conn,
	command *Command,
	firstResponse *datagram.DatagramRes,
	sendFileByteByByte func(net.Conn, string) error,
	readResponse func(net.Conn) (*datagram.DatagramRes, error),
	printResponse func(string, *datagram.DatagramRes),
) error {
	if command == nil {
		return fmt.Errorf("comando ADDFILE invalido")
	}

	if command.ID != datagram.ADDFILE {
		return fmt.Errorf("handler ADDFILE recebeu comando invalido: %d", command.ID)
	}

	if firstResponse == nil {
		return fmt.Errorf("primeira resposta do servidor ausente")
	}

	if firstResponse.StatusCode != datagram.STATUS_SUCCESS {
		return fmt.Errorf("ADDFILE recusado no primeiro status")
	}

	if err := sendFileByteByByte(conn, command.Filename); err != nil {
		return err
	}

	secondResponse, err := readResponse(conn)
	if err != nil {
		return err
	}

	printResponse("resposta 2", secondResponse)
	return nil
}
