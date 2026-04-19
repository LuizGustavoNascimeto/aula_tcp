package datagram

import "fmt"

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

type DatagramReq struct {
	MessageType  MessageType
	CommandID    CommandID
	FilenameSize int8
	Filename     string
	fileSize     uint32
}
type DatagramRes struct {
	MessageType MessageType
	CommandID   CommandID
	StatusCode  int8
}

func ParseReq(data []byte) (*DatagramReq, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("INVALID_DATAGRAM")
	}
	req := &DatagramReq{
		MessageType:  MessageType(data[0]),
		CommandID:    CommandID(data[1]),
		FilenameSize: int8(data[2]),
	}
	if req.MessageType != REQ_CODE {
		return nil, fmt.Errorf("INVALID_DATAGRAM_TYPE")
	}
	return req, nil
}

func (res *DatagramReq) String() string {
	return fmt.Sprintf("TYPE_%d CMD_%d FILENAME_SIZE_%d FILENAME_%s", res.MessageType, res.CommandID, res.FilenameSize, res.Filename)
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
	return res, nil
}

func CreateRes(commandID CommandID, statusCode int8) *DatagramRes {
	return &DatagramRes{
		MessageType: RES_CODE,
		CommandID:   commandID,
		StatusCode:  statusCode,
	}
}

func (res *DatagramReq) GetFileSize() uint32 {
	if res.CommandID != ADDFILE {
		return 0
	}
	return res.fileSize
}

func (res *DatagramReq) SetFileSize(size uint32) {
	if res.CommandID != ADDFILE {
		return
	}
	res.fileSize = size
}

// o datagram deve ter 3 bytes de cabeçalho
// E a té 256 bytes para o filename
