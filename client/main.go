package main

import (
	"cc_project/protocol/fstp"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
)

func main() {
	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	client, _ := fstp.NewFSTPClient(config)
	// body := fstp.IHaveProps{Files: []fstp.FileInfo{{Id: 1}}}
	// client.Request(fstp.IHaveRequest(body))
	fdata, _ := HashFile("/home/cralos/Uni/3Ano/CC/CC-Project/client_files/test_file.txt")
	fdata.OriginatorIP = client.Conn.LocalAddr().String()
	fdata.Name = "test_file.txt"
	
	client.Request(fstp.IHaveFileRequest(fstp.IHaveFileProps(*fdata)))
	// make_file_available("client_files/test_file.txt")
	// Send and receive data with the server
}
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

		// Print or process the hash of the chunk
		fmt.Printf("Hash of %d bytes: %x\n", n, chunkHash)

		// Update the hash of the entire file
		fileHash.Write(buffer[:n])

		// Append the chunkHash to the SegmentHashes field
		fileMetaData.SegmentHashes = append(fileMetaData.SegmentHashes, chunkHash)

		// Update the Length field
		fileMetaData.Length += int32(n)
	}

	// Get the final hash of the entire file
	fileChecksum := fstp.Hash(fileHash.Sum32()) // Convert to the Hash type
	fmt.Printf("Hash of the entire file: %x\n", fileChecksum)

	// Set the Hash field in the FileMetaData
	fileMetaData.Hash = fileChecksum
	return &fileMetaData, nil
}
