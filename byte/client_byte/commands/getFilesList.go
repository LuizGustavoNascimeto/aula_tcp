package commands

import (
	"client_byte/datagram"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// func HandleGetFilesListRes(
// 	conn net.Conn,
// 	resp *datagram.DatagramRes,
// 	readResponse func(net.Conn) (*datagram.DatagramRes, error),
// ) (*datagram.DatagramRes, error) {

// 	if resp == nil {
// 		return nil, fmt.Errorf("resposta GETFILESLIST vazia")
// 	}

// 	if resp.CommandID != datagram.GETFILESLIST {
// 		return nil, fmt.Errorf("resposta GETFILESLIST com comando invalido: %d", resp.CommandID)
// 	}

// 	if resp.StatusCode != datagram.STATUS_SUCCESS {
// 		return nil, fmt.Errorf("GETFILESLIST retornou erro")
// 	}
// 	newResponse, err := readResponse(conn)
// 	if err != nil {
// 		return newResponse, err
// 	}

// 	return newResponse, nil
// }

// HandleGetFilesListResponse executado após ler resposta GETFILESLIST
func HandleGetFilesListRes(conn net.Conn, res *datagram.DatagramRes) error {

	if res == nil {
		return fmt.Errorf("resposta GETFILESLIST vazia")
	}

	if res.CommandID != datagram.GETFILESLIST {
		return fmt.Errorf("resposta GETFILESLIST com comando invalido: %d", res.CommandID)
	}

	if res.StatusCode != datagram.STATUS_SUCCESS {
		return fmt.Errorf("GETFILESLIST retornou erro: %d", res.StatusCode)
	}

	numberOfFilesData := make([]byte, 2)
	if _, err := io.ReadFull(conn, numberOfFilesData); err != nil {
		return fmt.Errorf("falha ao ler número de arquivos: %w", err)
	}
	numberOfFiles := binary.BigEndian.Uint16(numberOfFilesData)
	res.SetNumberOfFiles(numberOfFiles)

	// Ler cada arquivo
	for i := 0; i < int(numberOfFiles); i++ {
		fileHeader := make([]byte, 1) // 1 byte para tamanho do nome
		if _, err := io.ReadFull(conn, fileHeader); err != nil {
			return fmt.Errorf("falha ao ler tamanho do nome do arquivo: %w", err)
		}
		filenameSize := uint8(fileHeader[0])

		filenameData := make([]byte, filenameSize)
		if _, err := io.ReadFull(conn, filenameData); err != nil {
			return fmt.Errorf("falha ao ler nome do arquivo: %w", err)
		}

		res.AddFile(datagram.FileStrcut{
			FilenameSize: filenameSize,
			Filename:     string(filenameData),
		})
	}

	return nil
}
