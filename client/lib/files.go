package lib

import (
	"cc_project/protocol/fstp"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
)

func HashFile(path string) (*fstp.FileMetaData, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return nil, err
	}
	defer file.Close()

	file_name := filepath.Base(path)

	// Create a buffer to read 128 bytes at a time
	buffer := make([]byte, 128)
	hasher := crc32.NewIEEE() // Use CRC32 hash

	// Hash the entire file
	fileHash := crc32.NewIEEE() // Use CRC32 hash

	// Create a FileMetaData instance
	fileMetaData := fstp.FileMetaData{
		// Initialize other fields as needed
		Name:          file_name,
		Length:        0, // To be updated
		OriginatorIP:  "PLACEHOLDER",
		SegmentHashes: []fstp.Hash{},
	}

	// Read and hash 128-byte chunks of the file
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break // Reached the end of the file
		}

		// Hash the chunk
		hasher.Reset()
		hasher.Write(buffer[:n])
		chunkHash := fstp.Hash(hasher.Sum32()) // Convert to the Hash type

		// Update the hash of the entire file
		fileHash.Write(buffer[:n])

		// Append the chunkHash to the SegmentHashes field
		fileMetaData.SegmentHashes = append(fileMetaData.SegmentHashes, chunkHash)

		// Update the Length field
		fileMetaData.Length += int32(n)
	}

	// Get the final hash of the entire file
	fileChecksum := fstp.FileHash(fileHash.Sum32()) // Convert to the Hash type

	// Set the Hash field in the FileMetaData
	fileMetaData.Hash = fileChecksum
	return &fileMetaData, nil
}

func MakeDirectoryAvailable(client *fstp.FSTPclient, directory string) {
	_, err := os.Stat(directory)

	visitFile := func(fp string, fi os.DirEntry, err error) error {
		if err != nil {
			return nil // but continue walking elsewhere
		}
		if fi.IsDir() {
			return nil // not a file. ignore.
		}

		// Call HashFile for the encountered file
		fdata, _ := HashFile(fp)
		if err != nil {
			fmt.Printf("Error hashing file %s: %v\n", fp, err)
		} else {
			fmt.Printf("File hashed: %s\n", fp)
			fmt.Printf("Hashed data: %+v\n", fdata)

			fdata.OriginatorIP = client.Conn.LocalAddr().String()
			fdata.Name = filepath.Base(fp)
			client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
		}

		return nil
	}

	if err == nil {
		// Directory exists
		fmt.Printf("Directory %s exists.\n", directory)
		err := filepath.WalkDir(directory, visitFile)
		if err != nil {
			fmt.Printf("error walking the path %v: %v\n", directory, err)
		}
	} else if os.IsNotExist(err) {
		// Directory does not exist
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
		} else {
			fmt.Printf("Directory %s created.\n", directory)
		}
	} else {
		// An error occurred (e.g., permission denied)
		fmt.Printf("Error checking directory: %v\n", err)
	}
}
