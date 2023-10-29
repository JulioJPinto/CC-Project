package db

import (
	"encoding/json"
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
