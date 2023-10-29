package db

import "net"

type DeviceData struct {
	Ip net.IP // Primary Key
}

type FileMetaData struct {
	Id   int64 // Primary Key
	Name string
}

type FileSegment struct {
	FirstByte int64
	Length    int64
	FileId    int64 // Foriegn Key refere um FileMetaData
}

func (s FileSegment) LastByte() int64 {
	return s.FirstByte + s.Length
}

