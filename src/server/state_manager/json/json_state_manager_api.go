package json

// import (
// 	"cc_project/helpers"
// 	"cc_project/protocol/fstp"
// 	"cc_project/server/state_manager"
// 	"fmt"
// 	"net"
// 	"strconv"
// 	"strings"
// )

// func (j *MapStateManager) RegisterDevice(data state_manager.DeviceData) error {
// 	devicesMap, ok := j.CachedDB["Devices"].(map[string]any)
// 	if ok {
// 		ip_string := data.Ip.String()
// 		if helpers.MapKeys[string](devicesMap).Contains(ip_string) {
// 			return nil
// 		}
// 		devicesMap[ip_string] = make([]string, 0)

// 	} else {
// 		e := state_manager.ErrBadSchema
// 		return e
// 	}

// 	return nil
// }

// func (db *MapStateManager) ResigerFile(file fstp.FileMetaData) {
// 	db.lock.Lock()
// 	defer db.lock.Unlock()

// 	filesMap, ok := db.CachedDB["Files"].(map[string]any)
// 	if ok {
// 		m := make(map[string]any)
// 		m["Name"] = file.Name
// 		filesMap[fmt.Sprint(file.HAsh)] = m
// 	}
// }

// func (db *MapStateManager) GetFileDataById(file_id int64) (*fstp.FileMetaData, error) {
// 	db.lock.RLock()
// 	defer db.lock.RUnlock()
// 	files, ok := db.CachedDB["Files"].(map[string]any)
// 	if !ok {

// 		return nil, state_manager.ErrBadSchema
// 	}

// 	file := new(fstp.FileMetaData)
// 	val, exists := files[fmt.Sprint(file_id)]
// 	if !exists {
// 		return nil, nil
// 	}
// 	m, ok := val.(map[string]any)
// 	if !ok {
// 		return nil, state_manager.ErrBadSchema
// 	}
// 	file.Name = fmt.Sprint(m["Name"])
// 	file.HAsh = file_id
// 	return file, nil
// }

// func (db *MapStateManager) _RegisterFileSegment(device_ip net.IP, segment fstp.FileSegment) error {

// 	db.lock.Lock()
// 	defer db.lock.Unlock()

// 	segmentsMap, ok := db.CachedDB["FileSegments"].(map[string]any)
// 	if !ok {
// 		return state_manager.ErrBadSchema
// 	}
// 	// fileData, err := db.GetFileDataById(segment.FileId)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// if fileData == nil {
// 	// 	return ErrFileDoesNotExist
// 	// }

// 	pk := fmt.Sprint(segment.FileId, "_", segment.FirstByte)
// 	curr, ok := segmentsMap[pk].(map[string]any)
// 	if !ok {
// 		segmentsMap[pk] = map[string]any{
// 			"ips": []net.IP{device_ip},
// 		}
// 	} else {
// 		slice := curr["ips"].([]any) // a bit iffy
// 		curr["ips"] = append(slice, device_ip)
// 		segmentsMap[pk] = curr
// 	}

// 	return nil
// }

// func (db *MapStateManager) RegisterFileSegment(device_ip net.IP, segment fstp.FileSegment) error {

// 	devicesMap, ok := db.CachedDB["Devices"].(map[string][]string)
// 	if ok {
// 		ip_string := device_ip.String()
// 		for k, v := range devicesMap {
// 			if k == ip_string {
// 				if helpers.SliceContains[string](v, fmt.Sprint(segment.FileId, "_", segment.FirstByte)) {
// 					return nil
// 				}
// 				v = append(v, fmt.Sprint(segment.FileId, "_", segment.FirstByte)) // alguem que checke
// 				devicesMap[k] = v
// 				return nil
// 			}
// 		}
// 		devicesMap[ip_string] = []string{fmt.Sprint(segment.FileId, "_", segment.FirstByte)}

// 		devicesMap[ip_string] = make([]string, 0)
// 		return nil
// 	} else {
// 		e := state_manager.ErrBadSchema
// 		return e
// 	}
// }

// // func (db *MapStateManager) BatchRegisterFileSegments(device_ip net.IP, segment []fstp.FileSegment) error {

// // 	devicesMap, ok := db.CachedDB["Devices"].(map[string][]string)
// // 	if ok {
// // 		ip_string := device_ip.String()
// // 		for k, v := range devicesMap {
// // 			if (k == ip_string) {
// // 				if helpers.SliceContains[string](v, fmt.Sprint(segment.FileId, "_", segment.FirstByte)) {
// //                     return nil
// //                 }
// // 				v = append(v,fmt.Sprint(segment.FileId, "_", segment.FirstByte)) // alguem que checke
// // 				devicesMap[k] = v
// // 				return nil
// // 			}
// // 		}
// // 		devicesMap[ip_string] = []string{fmt.Sprint(segment.FileId, "_", segment.FirstByte)}

// // 		devicesMap[ip_string] = make([]string, 0)
// // 		return nil
// // 	} else {
// // 		e := ErrBadSchema
// // 		return e
// // 	}
// // }

// func (db *MapStateManager) WhoHasFileSegment(segment fstp.FileSegment) ([]net.IP, error) {
// 	db.lock.RLock()
// 	defer db.lock.RUnlock()
// 	files_segments_map, ok := db.CachedDB["FileSegments"].(map[string]map[string]string)
// 	if !ok {
// 		return nil, state_manager.ErrBadSchema
// 	}
// 	ret := helpers.NewSet(func(a net.IP, b net.IP) bool { return a.Equal(b) })
// 	for key, element := range files_segments_map {
// 		fst_byte, _ := strconv.Atoi(strings.Split(key, "_")[1])
// 		if segment.FirstByte == int64(fst_byte) {
// 			ret.Add(net.ParseIP(element["ips"]))
// 		}
// 	}
// 	return ret.Slice(), nil
// }

// func (db *MapStateManager) RegisterFile(fstp.FileInfo) error {
// 	db.lock.Lock()
// 	defer db.lock.Unlock()

// 	_, ok := db.CachedDB["FileSegments"].(map[string]map[string]string)
// 	if !ok {
// 		return state_manager.ErrBadSchema
// 	}
// 	// files_segments_map

// 	return nil
// }
