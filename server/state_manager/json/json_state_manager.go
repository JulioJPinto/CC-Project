package json

// import (
// 	"cc_project/protocol/fstp"
// 	"cc_project/server/state_manager"
// 	"encoding/json"
// 	"os"
// 	"sync"
// )

// // writeToJSON writes a map to a JSON file
// func writeToJSON(data any, filePath string) error {
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}

// 	err = os.WriteFile(filePath, jsonData, 0644)
// 	return err
// }

// // readFromJSON reads a JSON file into a map
// func readFromJSON(filePath string) (map[string]interface{}, error) {
// 	jsonData, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var data map[string]interface{}
// 	err = json.Unmarshal(jsonData, &data)
// 	return data, err
// }

// // MapStateManager represents a JSON-based database
// type MapStateManager struct {
// 	lock     sync.RWMutex
// 	FilePath string
// 	CachedDB map[string]interface{}
// }

// func NewJSONDatabase(FilePath string) *MapStateManager {
// 	Map := make(map[string]any)
// 	Map["Devices"] = make(map[string]state_manager.DeviceData)
// 	Map["Files"] = make(map[string]fstp.FileMetaData)
// 	Map["FileSegments"] = make(map[string]fstp.FileSegment)
// 	j := &MapStateManager{
// 		sync.RWMutex{},
// 		FilePath,
// 		Map,
// 	}

// 	return j
// }

// // Connect opens the JSON database connection
// func (db *MapStateManager) Connect() error {
// 	db.lock.Lock()
// 	defer db.lock.Unlock()
// 	_, err := os.Stat(db.FilePath)

// 	if os.IsNotExist(err) {
// 		// If the file does not exist, create an empty one
// 		file, createErr := os.Create(db.FilePath)
// 		if createErr != nil {
// 			return createErr
// 		}
// 		writeToJSON(db.CachedDB, db.FilePath)
// 		defer file.Close()
// 	}
// 	data, err := readFromJSON(db.FilePath)
// 	if err != nil {
// 		return err
// 	}
// 	db.CachedDB = data
// 	// Additional initialization, such as initializing data structures or handling other setup tasks

// 	return nil
// }

// // Close closes the JSON database connection
// func (db *MapStateManager) Close() error {
// 	db.lock.Lock()
// 	defer db.lock.Unlock()
// 	return writeToJSON(db.CachedDB, db.FilePath)
// }

// func (db *MapStateManager) Write() error {
// 	db.lock.Lock()
// 	defer db.lock.Unlock()
// 	return writeToJSON(db.CachedDB, db.FilePath)
// }
