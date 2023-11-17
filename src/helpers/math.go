package helpers

func TrunkI64(n int64)int32{
	return int32(n & 0xFFFFFFFF)
}