package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false // não existe
		}
		// outro erro (permissão, etc.)
		return false
	}
	return info.IsDir() // garante que é um diretório, não um arquivo
}

func ListFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func SendFile(conn net.Conn, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("nao foi possivel abrir arquivo '%s': %w", filename, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		b, err := reader.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("falha ao ler arquivo '%s': %w", filename, err)
		}

		if _, err := conn.Write([]byte{b}); err != nil {
			return fmt.Errorf("falha ao enviar byte do arquivo '%s': %w", filename, err)
		}
	}

	fmt.Printf("envio byte a byte finalizado: %s\n", filename)
	return nil
}
