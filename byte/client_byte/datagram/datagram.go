package datagram

import (
	"encoding/binary"
	"fmt"
)

type DatagramReq struct {
	MessageType  MessageType
	CommandID    CommandID
	FilenameSize int8
	FileSize     uint32
	Filename     string
}
type DatagramRes struct {
	MessageType MessageType
	CommandID   CommandID
	StatusCode  int8
}

type MessageType int8
type CommandID int8

const (
	REQ_CODE        MessageType = 1
	RES_CODE        MessageType = 2
	INVALID_COMMAND CommandID   = 0
	ADDFILE         CommandID   = 1
	DELETE          CommandID   = 2
	GETFILESLIST    CommandID   = 3
	GETFILE         CommandID   = 4
	STATUS_SUCCESS  int8        = 0
	STATUS_ERROR    int8        = 1
)

func CreateReq(commandID CommandID, filename string, fileSize uint32) ([]byte, error) {
	if len(filename) > 255 {
		return nil, fmt.Errorf("FILENAME_TOO_LONG")
	}

	payloadLen := 3 + len(filename)
	if commandID == ADDFILE {
		payloadLen += 4
	}

	payload := make([]byte, 0, payloadLen)
	payload = append(payload, byte(REQ_CODE), byte(commandID), byte(len(filename)))
	payload = append(payload, []byte(filename)...)
	if commandID == ADDFILE {
		fileSizeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(fileSizeBytes, fileSize)
		payload = append(payload, fileSizeBytes...)
	}
	return payload, nil
}

func ParseRes(data []byte) (*DatagramRes, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("INVALID_DATAGRAM")
	}
	res := &DatagramRes{
		MessageType: MessageType(data[0]),
		CommandID:   CommandID(data[1]),
		StatusCode:  int8(data[2]),
	}
	if res.MessageType != RES_CODE {
		return nil, fmt.Errorf("INVALID_DATAGRAM_TYPE")
	}
	return res, nil
}
