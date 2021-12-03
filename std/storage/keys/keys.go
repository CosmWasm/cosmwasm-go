package keys

func StringPrimaryKey(s string) []byte {
	return []byte(s)
}

func BytesPrimaryKey(b []byte) []byte {
	return b
}

func BoolPrimaryKey(b bool) []byte {
	switch b {
	case true:
		return []byte{0x1}
	default:
		return []byte{0x0}
	}
}

func Int8PrimaryKey(i int8) []byte {
	return []byte{uint8(i)}
}

func Int16PrimaryKey(i int16) []byte {
	return Uint16PrimaryKey(uint16(i))
}

func Int32PrimaryKey(i int32) []byte {
	return Uint32PrimaryKey(uint32(i))
}

func Int64PrimaryKey(i int64) []byte {
	return Uint64PrimaryKey(uint64(i))
}

func Uint8PrimaryKey(i uint8) []byte {
	return []byte{i}
}

func Uint16PrimaryKey(i uint16) []byte {
	b := make([]byte, 2)
	b[0] = byte(i >> 8)
	b[1] = byte(i)
	return b
}

func Uint32PrimaryKey(i uint32) []byte {
	b := make([]byte, 4)

	b[0] = byte(i >> 24)
	b[1] = byte(i >> 16)
	b[2] = byte(i >> 8)
	b[3] = byte(i)

	return b
}

func Uint64PrimaryKey(i uint64) []byte {
	b := make([]byte, 8)

	b[0] = byte(i >> 56)
	b[1] = byte(i >> 48)
	b[2] = byte(i >> 40)
	b[3] = byte(i >> 32)
	b[4] = byte(i >> 24)
	b[5] = byte(i >> 16)
	b[6] = byte(i >> 8)
	b[7] = byte(i)

	return b
}
