package db

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// writeToJSON writes a map to a JSON file
func writeToJSON(data any, filePath string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	return err
}

// readFromJSON reads a JSON file into a map
func readFromJSON(filePath string) (map[string]interface{}, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	return data, err
}

// JSONDatabase represents a JSON-based database
type JSONDatabase struct {
	FilePath string
	CachedDB map[string]interface{}
}

func NewJSONDatabase(FilePath string) *JSONDatabase {
	Map := make(map[string]any)
	Map["Devices"] = make(map[string]DeviceData)
	Map["Files"] = make(map[string]FileMetaData)
	Map["FileSegments"] = make(map[string]FileSegment)
	Map["DevicesFileSegments"] = make(map[string]DevicesFileSegments)
	j := &JSONDatabase{
		FilePath,
		Map,
	}

	return j
}

// Connect opens the JSON database connection
func (j *JSONDatabase) Connect() error {
	_, err := os.Stat(j.FilePath)

	if os.IsNotExist(err) {
		// If the file does not exist, create an empty one
		file, createErr := os.Create(j.FilePath)
		if createErr != nil {
			return createErr
		}
		writeToJSON(j.CachedDB, j.FilePath)
		defer file.Close()
	}
	data, err := readFromJSON(j.FilePath)
	if err != nil {
		return err
	}
	j.CachedDB = data
	// Additional initialization, such as initializing data structures or handling other setup tasks

	return nil
}

// Close closes the JSON database connection
func (j *JSONDatabase) Close() error {
	return writeToJSON(j.CachedDB, j.FilePath)
}

func (j *JSONDatabase) RegisterDevice(data DeviceData) error {
	devicesMap, ok := j.CachedDB["Devices"].(map[string]any)
	if ok {
		devicesMap[data.Ip.String()] = make([]string, 0)
	} else {
		e := ErrBadSchema
		return e
	}

	return nil
}

func (j *JSONDatabase) ResigerFile(file FileMetaData) {
	filesMap, ok := j.CachedDB["Files"].(map[string]any)
	if ok {
		m := make(map[string]any)
		m["Name"] = file.Name
		filesMap[fmt.Sprint(file.Id)] = m
	}
}

func (j *JSONDatabase) GetFileData(file_id int64) (*FileMetaData, error) {
	files, ok := j.CachedDB["Files"].(map[string]any)
	if !ok {

		return nil, ErrBadSchema
	}

	file := new(FileMetaData)
	val, exists := files[fmt.Sprint(file_id)]
	if !exists {
		return nil, nil
	}
	m, ok := val.(map[string]any)
	if !ok {
		return nil, ErrBadSchema
	}
	file.Name = fmt.Sprint(m["Name"])
	file.Id = file_id
	return file, nil
}

/*
Map["FileSegments"] = make(map[string]FileSegment)
Map["DevicesFileSegments"] = make(map[string]DevicesFileSegments)
*/
func (j *JSONDatabase) RegisterFileSegment(device_ip net.IP, segment FileSegment) error {
	segmentsMap, ok := j.CachedDB["FileSegments"].(map[string]any)
	if !ok {
		return ErrBadSchema
	}
	fileData, err := j.GetFileData(segment.FileId)
	if err != nil {
		return err
	}
	if fileData == nil {
		return ErrFileDoesNotExist
	}
	pk := fmt.Sprint(segment.FileId, "_", segment.FirstByte)
	segmentsMap[pk] = map[string]string{
		"length": fmt.Sprint(segment.Length),
	}
	// segmentsMap[device_ip.String()] = make([]string, 0)
	return nil
}
