package fstp

import (
	"cc_project/protocol"
	"fmt"
	"os"
	"path/filepath"
)

func HashFile(path string) (*protocol.FileMetaData, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return nil, err
	}
	defer file.Close()

	file_name := filepath.Base(path)

	return protocol.HashFile(file, file_name)
}
