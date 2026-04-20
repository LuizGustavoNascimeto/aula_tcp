package datagram

import (
	"client_byte/utils"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type FileStrcut struct {
	FilenameSize uint8
	Filename     string
}

type DatagramReq struct {
	MessageType  MessageType
	CommandID    CommandID
	FilenameSize int8
	FileSize     uint32
	Filename     string
}
type DatagramRes struct {
	MessageType   MessageType
	CommandID     CommandID
	StatusCode    int8
	numberOfFiles uint16
	files         []FileStrcut
	fileSize      uint32
	fileBytes     []byte
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
	STATUS_SUCCESS  int8        = 1
	STATUS_ERROR    int8        = 2
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

	switch res.CommandID {
	case ADDFILE, DELETE, GETFILE:
		if res.StatusCode != STATUS_SUCCESS {
			return res, nil
		}
	case GETFILESLIST:
		if res.StatusCode != STATUS_SUCCESS {
			return res, nil
		}

	default:
		return nil, fmt.Errorf("UNKNOWN_COMMAND_ID")
	}
	return res, nil
}

func (res *DatagramReq) GetFileSize() uint32 {
	if res.CommandID != ADDFILE {
		return 0
	}
	return res.FileSize
}

func (res *DatagramRes) SetFiles(files []FileStrcut) {
	if res.CommandID != GETFILESLIST {
		return
	}
	res.numberOfFiles = uint16(len(files))
	res.files = files
}

func (res *DatagramRes) GetNumberOfFiles() uint16 {
	if res.CommandID != GETFILESLIST {
		return 0
	}
	return res.numberOfFiles
}

func (res *DatagramRes) SetNumberOfFiles(count uint16) {
	res.numberOfFiles = count
}

func (res *DatagramRes) AddFile(file FileStrcut) {
	res.files = append(res.files, file)
}

func (res *DatagramRes) SetFilePayload(size uint32, bytes []byte) {
	if res.CommandID != GETFILE {
		return
	}
	res.fileSize = size
	res.fileBytes = bytes
}

func (res *DatagramRes) GetFileSize() uint32 {
	if res.CommandID != GETFILE {
		return 0
	}
	return res.fileSize
}

func (res *DatagramRes) GetFileBytes() []byte {
	if res.CommandID != GETFILE {
		return nil
	}
	return res.fileBytes
}

func (res *DatagramRes) GetFilesList() []FileStrcut {
	if res.CommandID != GETFILESLIST {
		return nil
	}
	return res.files
}

func ParseReq(msg string) (*DatagramReq, error) {
	fields := strings.Fields(msg)
	if len(fields) == 0 {
		return nil, errors.New("INVALID_COMMAND")
	}

	commandID, needsFilename, err := readCommandID(fields[0])
	if err != nil {
		return nil, err
	}

	filename := ""
	if needsFilename {
		if len(fields) < 2 {
			return nil, errors.New("FILENAME_REQUIRED")
		}
		filename = strings.Join(fields[1:], " ")

		if _, err := utils.ReadFilenameSize(filename); err != nil {
			return nil, err
		}
	}
	fileSize := uint32(0)
	if commandID == ADDFILE {
		fileSize, err = utils.GetFileSize(filename)
		if err != nil {
			return nil, err
		}
	}

	return &DatagramReq{
		CommandID:    commandID,
		FilenameSize: int8(len(filename)),
		Filename:     filename,
		FileSize:     fileSize,
	}, nil
}

func readCommandID(msg string) (CommandID, bool, error) {
	switch strings.ToUpper(msg) {
	case "ADDFILE":
		return ADDFILE, true, nil
	case "DELETE":
		return DELETE, true, nil
	case "GETFILESLIST":
		return GETFILESLIST, false, nil
	case "GETFILE":
		return GETFILE, true, nil
	default:
		return INVALID_COMMAND, false, errors.New("INVALID_COMMAND")
	}
}
