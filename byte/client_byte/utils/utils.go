package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func resolveClientFilePath(filename string) string {
	clean := strings.TrimSpace(filename)
	if clean == "" {
		return filename
	}

	if filepath.IsAbs(clean) || strings.Contains(clean, "/") || strings.Contains(clean, "\\") {
		return clean
	}

	return filepath.Join("files", clean)
}

func SendFile(conn net.Conn, filename string) error {
	filePath := resolveClientFilePath(filename)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("nao foi possivel abrir arquivo '%s': %w", filePath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		b, err := reader.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("falha ao ler arquivo '%s': %w", filePath, err)
		}

		if _, err := conn.Write([]byte{b}); err != nil {
			return fmt.Errorf("falha ao enviar byte do arquivo '%s': %w", filePath, err)
		}
	}

	fmt.Printf("envio byte a byte finalizado: %s\n", filePath)
	return nil
}

func ReadFilenameSize(filename string) (int8, error) {
	if len(filename) > 255 {
		return 0, errors.New("FILENAME_TOO_LONG")
	}
	return int8(len(filename)), nil
}

func GetFileSize(filename string) (uint32, error) {
	filePath := resolveClientFilePath(filename)
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("nao foi possivel obter informacoes do arquivo '%s': %w", filePath, err)
	}
	return uint32(info.Size()), nil
}

func SaveDownloadedFile(filename string, fileData []byte) (string, error) {
	clean := filepath.Base(strings.TrimSpace(filename))
	if clean == "" || clean == "." || clean == string(filepath.Separator) {
		return "", errors.New("FILENAME_INVALID")
	}

	targetDir := "files"
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("nao foi possivel criar pasta '%s': %w", targetDir, err)
	}

	targetPath := filepath.Join(targetDir, clean)
	if err := os.WriteFile(targetPath, fileData, 0644); err != nil {
		return "", fmt.Errorf("nao foi possivel salvar arquivo '%s': %w", targetPath, err)
	}

	return targetPath, nil
}
