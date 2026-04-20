package commands

import (
	"math"
	"os"
	"path/filepath"
	"tcp_server_byte/src/datagram"
)

const downloadDir = "files"

func GETFILE(filename string) ([]byte, int8, error) {
	safeName := filepath.Base(filename)
	targetPath := filepath.Join(downloadDir, safeName)

	fileBytes, err := os.ReadFile(targetPath)
	if err != nil {
		return nil, datagram.STATUS_ERROR, err
	}
	if len(fileBytes) > math.MaxUint32 {
		return nil, datagram.STATUS_ERROR, os.ErrInvalid
	}

	return fileBytes, datagram.STATUS_SUCCESS, nil
}
