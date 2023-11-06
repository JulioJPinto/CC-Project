package fstp

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
)

func HashFile(path string) (*FileMetaData, error) {
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
	fileMetaData := FileMetaData{
		// Initialize other fields as needed
		Name:          file_name,
		Length:        0, // To be updated
		OriginatorIP:  "PLACEHOLDER",
		SegmentHashes: []Hash{},
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
		chunkHash := Hash(hasher.Sum32()) // Convert to the Hash type

		// Update the hash of the entire file
		fileHash.Write(buffer[:n])

		// Append the chunkHash to the SegmentHashes field
		fileMetaData.SegmentHashes = append(fileMetaData.SegmentHashes, chunkHash)

		// Update the Length field
		fileMetaData.Length += int32(n)
	}

	// Get the final hash of the entire file
	fileChecksum := FileHash(fileHash.Sum32()) // Convert to the Hash type

	// Set the Hash field in the FileMetaData
	fileMetaData.Hash = fileChecksum
	return &fileMetaData, nil
}