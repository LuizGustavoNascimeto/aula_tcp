package client

import (
	"client_byte/datagram"
	"fmt"
	"time"
)

const originClient = "cliente"

func AppLog(origin string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[%s][%s] %s\n", timestamp, origin, msg)
}

func ClientLog(format string, args ...interface{}) {
	AppLog(originClient, format, args...)
}

func PrintResponse(res *datagram.DatagramRes) {
	if res == nil {
		AppLog("servidor", "resposta vazia")
		return
	}

	switch res.CommandID {
	case datagram.ADDFILE:
		AppLog("servidor", "cmd=%d status=%d", res.CommandID, res.StatusCode)
	case datagram.GETFILESLIST:
		PrintFilesListResponse(res)
	default:
		AppLog("servidor", "cmd=%d status=%d", res.CommandID, res.StatusCode)
	}
}

func PrintFilesListResponse(res *datagram.DatagramRes) {
	AppLog("servidor", "cmd=%d status=%d numberOfFiles=%d", res.CommandID, res.StatusCode, res.GetNumberOfFiles())
	for i, file := range res.GetFilesList() {
		AppLog("servidor", "file %d: %s", i+1, file.Filename)
	}
}
