package db

import (
	"cc_project/helpers"
	"fmt"
	"net"
	"strconv"
	"strings"
)

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

func (j *JSONDatabase) GetFileDataById(file_id int64) (*FileMetaData, error) {
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

func (j *JSONDatabase) RegisterFileSegment(device_ip net.IP, segment FileSegment) error {
	segmentsMap, ok := j.CachedDB["FileSegments"].(map[string]any)
	if !ok {
		return ErrBadSchema
	}
	fileData, err := j.GetFileDataById(segment.FileId)
	if err != nil {
		return err
	}
	if fileData == nil {
		return ErrFileDoesNotExist
	}
	pk := fmt.Sprint(segment.FileId, "_", segment.FirstByte)
	curr, ok := segmentsMap[pk].(map[string]any)
	if !ok {
		segmentsMap[pk] = map[string]any{
			"length": segment.Length,
			"ips":    []net.IP{device_ip},
		}
	} else {
		slice := curr["ips"].([]any) // a bit iffy
		curr["ips"] = append(slice, device_ip)
		segmentsMap[pk] = curr
	}

	return nil
}

func (j *JSONDatabase) WhoHasFileSegment(segment FileSegment) ([]net.IP, error) {
	files_segments_map, ok := j.CachedDB["FileSegments"].(map[string]map[string]string)
	if !ok {
		return nil, ErrBadSchema
	}
	ret := helpers.NewSet(func(a net.IP, b net.IP) bool { return a.Equal(b) })
	for key, element := range files_segments_map {
		fst_byte, _ := strconv.Atoi(strings.Split(key, "_")[1])
		length, _ := strconv.Atoi(element["length"])
		last_byte := fst_byte + length
		if segment.FirstByte >= int64(fst_byte) && segment.LastByte() <= int64(last_byte) {
			ret.Add(net.ParseIP(element["ips"]))
		}
	}
	return ret.Slice(), nil
}
