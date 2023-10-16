package db

import "net"

type DeviceData struct {
	Ip net.IP
}

type FileMetaData struct{
	name string
	id int64
}

type FileSegment struct{
	firstByte int64
	length int64
}
