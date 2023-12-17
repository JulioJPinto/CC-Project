package protocol

import (
	"hash/crc32"
	"os"
)

func HashSegment(buffer []byte, n int) Hash {
	hasher := crc32.NewIEEE() // Use CRC32 hash

	hasher.Reset()
	hasher.Write(buffer[:n])
	return Hash(hasher.Sum32()) // Convert to the Hash type
}

func Hashing(file *os.File, file_name string) (*FileMetaData, error) {
	// Create a buffer to read 128 bytes at a time
	buffer := make([]byte, SegmentMaxLength)

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
		chunkHash := HashSegment(buffer, n)
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
