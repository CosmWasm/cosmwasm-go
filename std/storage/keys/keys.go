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

// Int8PrimaryKey returns a byte int8 representation
// such as -128 <= i <= 127
func Int8PrimaryKey(i int8) []byte {
	b := uint8(i)
	b ^= 0x80
	return []byte{b}
}

// Int16PrimaryKey returns a byte int16 representation
// such as -32,768 <= i <= 32,767
func Int16PrimaryKey(i int16) []byte {
	bytes := Uint16PrimaryKey(uint16(i))
	bytes[0] ^= 0x80
	return bytes
}

// Int32PrimaryKey returns a byte int32 representation
// such as -2,147,483,648 <= i <= 2,147,483,647
func Int32PrimaryKey(i int32) []byte {
	bytes := Uint32PrimaryKey(uint32(i))
	bytes[0] ^= 0x80
	return bytes
}

// Int64PrimaryKey returns a byte int64 representation
// such as -9,223,372,036,854,775,808 <= i <= 9,223,372,036,854,775,807
func Int64PrimaryKey(i int64) []byte {
	bytes := Uint64PrimaryKey(uint64(i))
	bytes[0] ^= 0x80
	return bytes
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
