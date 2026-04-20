package commands

import (
	"bufio"
	"tcp_server_byte/src/datagram"
	"tcp_server_byte/src/utils"
)

// essa função faz updaload de um arquivo para o servidor
func GETFILESLIST(reader *bufio.Reader) ([]datagram.FileStrcut, error) {
	files, err := utils.ListFiles("files")
	if err != nil {
		return nil, err
	}
	var res []datagram.FileStrcut
	for _, file := range files {
		res = append(res, datagram.FileStrcut{
			Filename:     file,
			FilenameSize: uint8(len(file)),
		})

	}
	return res, nil
}
