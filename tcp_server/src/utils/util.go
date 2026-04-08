package utils

import "os"

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
