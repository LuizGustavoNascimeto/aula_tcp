package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"

	"client_byte/commands"
	"client_byte/datagram"
)

func processCommand(conn net.Conn, rawInput string) error {
	command, err := commands.Parse(rawInput)
	if err != nil {
		return err
	}

	if err := sendDatagram(conn, command); err != nil {
		return err
	}

	resp, err := readResponse(conn)
	if err != nil {
		return err
	}

	printResponse("resposta 1", resp)

	switch command.ID {
	case datagram.ADDFILE:
		fmt.Printf("comando ADDFILE enviado, processando upload...\n")
		if err := commands.HandleADDFILE(conn, command, resp, sendFileByteByByte, readResponse, printResponse); err != nil {
			return err
		}
	default:
		fmt.Printf("comando %d finalizado\n", command.ID)
	}

	return nil

}

func sendDatagram(conn net.Conn, command *commands.Command) error {
	fileSize := uint32(0)
	if command.ID == datagram.ADDFILE {
		size, err := fileSizeFromPath(command.Filename)
		if err != nil {
			return err
		}
		fileSize = size
	}

	payload, err := datagram.CreateReq(command.ID, command.Filename, fileSize)
	if err != nil {
		return err
	}

	_, err = conn.Write(payload)
	return err
}

func fileSizeFromPath(filename string) (uint32, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, fmt.Errorf("nao foi possivel ler tamanho do arquivo '%s': %w", filename, err)
	}

	if info.IsDir() {
		return 0, fmt.Errorf("'%s' eh diretorio, esperado arquivo", filename)
	}

	if info.Size() > math.MaxUint32 {
		return 0, fmt.Errorf("arquivo '%s' excede limite de 4 bytes para tamanho", filename)
	}

	return uint32(info.Size()), nil
}

func readResponse(conn net.Conn) (*datagram.DatagramRes, error) {
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

func printResponse(label string, res *datagram.DatagramRes) {
	fmt.Printf("server (%s): cmd=%d status=%d\n", label, res.CommandID, res.StatusCode)
}

func sendFileByteByByte(conn net.Conn, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("nao foi possivel abrir arquivo '%s': %w", filename, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		b, err := reader.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("falha ao ler arquivo '%s': %w", filename, err)
		}

		if _, err := conn.Write([]byte{b}); err != nil {
			return fmt.Errorf("falha ao enviar byte do arquivo '%s': %w", filename, err)
		}
	}

	fmt.Printf("envio byte a byte finalizado: %s\n", filename)
	return nil
}
