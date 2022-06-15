package gojce

func JceTypeIsInt8(jceType uint8) bool {
	return jceType == 0 || jceType == 12
}

func JceTypeIsInt16(jceType uint8) bool {
	return jceType == 0 || jceType == 1 || jceType == 12
}

func JceTypeIsInt32(jceType uint8) bool {
	return jceType == 0 || jceType == 1 || jceType == 2 || jceType == 12
}

func JceTypeIsInt64(jceType uint8) bool {
	return jceType == 0 || jceType == 1 || jceType == 2 || jceType == 3 || jceType == 12
}

func JceTypeIsFloat32(jceType uint8) bool {
	return jceType == 4
}

func JceTypeIsFloat64(jceType uint8) bool {
	return jceType == 5
}

func JceTypeIsString(jceType uint8) bool {
	return jceType == 6 || jceType == 7
}

func JceTypeIsMap(jceType uint8) bool {
	return jceType == 8
}

func JceTypeIsList(jceType uint8) bool {
	return jceType == 9
}

func JceTypeIsStruct(jceType uint8) bool {
	return jceType == 10
}

func JceTypeIsBytes(jceType uint8) bool {
	return jceType == 13
}
