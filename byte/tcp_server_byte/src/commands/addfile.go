package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"tcp_server_byte/src/datagram"
)

const uploadDir = "files"

// essa função faz updaload de um arquivo para o servidor
func ADDFILE(reader *bufio.Reader, filename string, filesize uint32) (int8, error) {
	file := readFileFromConn(reader, filesize)
	fmt.Println("arquivo lido")

	safeName := filepath.Base(filename)
	targetPath := filepath.Join(uploadDir, safeName)

	err := os.WriteFile(targetPath, file, 0644)
	if err != nil {
		return datagram.STATUS_ERROR, err
	}
	return datagram.STATUS_SUCCESS, nil
}

// Lê o arquivo byte a byte e retorna em uma variável

func readFileFromConn(reader *bufio.Reader, size uint32) []byte {
	var data []byte

	for uint32(len(data)) < size {
		b, err := reader.ReadByte()
		if err != nil {
			// outro erro (conexão resetada, timeout, etc.)
			break
		}
		data = append(data, b)
		fmt.Println("byte lido")
	}
	return data
}
