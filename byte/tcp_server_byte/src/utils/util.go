package utils

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

var basePath = "users_files"

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

// verifica que o caminho pode ser acessado pelo usuário
func IsValidDir(root string, newDir string) bool {
	cleanRoot := path.Clean("/" + strings.TrimPrefix(root, "/"))
	cleanNewDir := path.Clean("/" + strings.TrimPrefix(newDir, "/"))

	if cleanNewDir != cleanRoot && !strings.HasPrefix(cleanNewDir, cleanRoot+"/") {
		return false
	}

	relativeDir := filepath.FromSlash(strings.TrimPrefix(cleanNewDir, "/"))
	targetDir := filepath.Join(basePath, relativeDir)

	return DirExists(targetDir)
}

func ListFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(basePath, dir))
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

func ListDirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(basePath, dir))
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, entry := range entries {

		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs, nil
}
