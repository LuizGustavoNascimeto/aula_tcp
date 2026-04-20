package datagram

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

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

type FileStrcut struct {
	FilenameSize uint8
	Filename     string
}
type DatagramReq struct {
	MessageType  MessageType
	CommandID    CommandID
	FilenameSize uint8
	Filename     string
	fileSize     uint32
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

func ParseReq(data []byte) (*DatagramReq, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("INVALID_DATAGRAM")
	}
	req := &DatagramReq{
		MessageType:  MessageType(data[0]),
		CommandID:    CommandID(data[1]),
		FilenameSize: uint8(data[2]),
	}
	return req, nil
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

func (res *DatagramRes) SetFiles(files []FileStrcut) {
	if res.CommandID != GETFILESLIST {
		return
	}
	res.numberOfFiles = uint16(len(files))
	res.files = files

}

func (res *DatagramRes) SetFileBytes(file []byte) {
	if res.CommandID != GETFILE {
		return
	}
	res.fileSize = uint32(len(file))
	res.fileBytes = file
}

func (res *DatagramReq) String() string {
	return fmt.Sprintf("TYPE_%d CMD_%d FILENAME_SIZE_%d FILENAME_%s", res.MessageType, res.CommandID, res.FilenameSize, res.Filename)
}

func (req *DatagramReq) HandleFilename(reader *bufio.Reader) error {
	// Lê o nome do arquivo, se houver
	if req.FilenameSize <= 0 {
		return nil
	}
	filenameBytes := make([]byte, int(req.FilenameSize))
	if _, err := io.ReadFull(reader, filenameBytes); err != nil {
		return err
	}

	req.Filename = string(filenameBytes)

	if req.CommandID == ADDFILE {
		// Para ADDFILE, também lemos os próximos 4 bytes para o tamanho do arquivo
		fileSizeBytes := make([]byte, 4)
		if _, err := io.ReadFull(reader, fileSizeBytes); err != nil {
			return err
		}
		req.fileSize = binary.BigEndian.Uint32(fileSizeBytes)
	}

	return nil
}

func (res *DatagramRes) ToBytes() []byte {
	payload := make([]byte, 3)
	payload[0] = byte(res.MessageType)
	payload[1] = byte(res.CommandID)
	payload[2] = byte(res.StatusCode)
	if res.CommandID == GETFILESLIST && res.StatusCode == STATUS_SUCCESS {
		// Serializa apenas nomes válidos no protocolo: tamanho entre 1 e 255 bytes.
		validFiles := make([][]byte, 0, len(res.files))
		for _, file := range res.files {
			nameBytes := []byte(file.Filename)
			if len(nameBytes) < 1 || len(nameBytes) > 255 {
				continue
			}
			validFiles = append(validFiles, nameBytes)
		}

		numberOfFilesBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(numberOfFilesBytes, uint16(len(validFiles)))
		payload = append(payload, numberOfFilesBytes...)

		for _, filenameBytes := range validFiles {
			payload = append(payload, byte(len(filenameBytes)))
			payload = append(payload, filenameBytes...)
		}
	}
	if res.CommandID == GETFILE && res.StatusCode == STATUS_SUCCESS {
		fileSizeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(fileSizeBytes, res.fileSize)
		payload = append(payload, fileSizeBytes...)
		payload = append(payload, res.fileBytes...)
	}
	return payload
}
