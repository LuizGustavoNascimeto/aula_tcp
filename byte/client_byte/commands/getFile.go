package commands

import (
	"client_byte/datagram"
	"client_byte/utils"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func HandleGetFileRes(conn net.Conn, filename string, res *datagram.DatagramRes) error {
	if res == nil {
		return fmt.Errorf("resposta GETFILE vazia")
	}
	if res.CommandID != datagram.GETFILE {
		return fmt.Errorf("resposta GETFILE com comando invalido: %d", res.CommandID)
	}
	if res.StatusCode != datagram.STATUS_SUCCESS {
		return fmt.Errorf("GETFILE retornou erro: %d", res.StatusCode)
	}

	fileSizeData := make([]byte, 4)
	if _, err := io.ReadFull(conn, fileSizeData); err != nil {
		return fmt.Errorf("falha ao ler tamanho do arquivo: %w", err)
	}

	fileSize := binary.BigEndian.Uint32(fileSizeData)
	fileData := make([]byte, fileSize)
	if _, err := io.ReadFull(conn, fileData); err != nil {
		return fmt.Errorf("falha ao ler bytes do arquivo: %w", err)
	}

	res.SetFilePayload(fileSize, fileData)
	if _, err := utils.SaveDownloadedFile(filename, fileData); err != nil {
		return err
	}

	return nil
}
