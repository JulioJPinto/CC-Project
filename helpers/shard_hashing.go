package helpers

type HashFn[K comparable] func(k K) uint64

func HashString(s string) uint64 {
	var hash uint64 = 5381

	for _, c := range s {
		hash = ((hash << 5) + hash) + uint64(c)
	}

	return hash
}

func HashUint64(u uint64) uint64 {
	u ^= u >> 33
	u *= 0xff51afd7ed558ccd
	u ^= u >> 33
	u *= 0xc4ceb9fe1a85ec53
	u ^= u >> 33
	return u
}

func HashUint32(u uint32) uint64 {
	return HashUint64(uint64(u))
}


