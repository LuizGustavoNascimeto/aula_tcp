package commands

import (
	"errors"
	"strings"

	"client_byte/datagram"
)

type Command struct {
	ID       datagram.CommandID
	Raw      string
	Filename string
}

func Parse(msg string) (*Command, error) {
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

		if _, err := readFilenameSize(filename); err != nil {
			return nil, err
		}
	}

	return &Command{
		ID:       commandID,
		Raw:      msg,
		Filename: filename,
	}, nil
}

func readCommandID(msg string) (datagram.CommandID, bool, error) {
	switch strings.ToUpper(msg) {
	case "ADDFILE":
		return datagram.ADDFILE, true, nil
	case "DELETE":
		return datagram.DELETE, true, nil
	case "GETFILESLIST":
		return datagram.GETFILESLIST, false, nil
	case "GETFILE":
		return datagram.GETFILE, true, nil
	default:
		return datagram.INVALID_COMMAND, false, errors.New("INVALID_COMMAND")
	}
}

func readFilenameSize(filename string) (int8, error) {
	if len(filename) > 255 {
		return 0, errors.New("FILENAME_TOO_LONG")
	}
	return int8(len(filename)), nil
}
