package commands

import (
	"os"
	"path/filepath"
	"tcp_server_byte/src/datagram"
)

const filesDir = "files"

func DELETE(filename string) (int8, error) {
	safeName := filepath.Base(filename)
	targetPath := filepath.Join(filesDir, safeName)

	if err := os.Remove(targetPath); err != nil {
		return datagram.STATUS_ERROR, err
	}

	return datagram.STATUS_SUCCESS, nil
}
