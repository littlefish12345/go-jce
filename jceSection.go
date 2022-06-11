package gojce

const (
	INT8              = "INT8"
	BYTE              = "BYTE"
	INT16             = "INT16"
	INT32             = "INT32"
	INT64             = "INT64"
	FLOAT32           = "FLOAT32"
	FLOAT64           = "FLOAT64"
	STRING            = "STRING"
	MAP               = "MAP" //别用, 只会在section解码的时候出现
	MAPStrStr         = "MAPStrStr"
	MAPStrBytes       = "MAPStrBytes"
	MAPStrMAPStrBytes = "MAPStrMAPStrBytes"
	LIST              = "LIST" //别用, 只会在section解码的时候出现
	LISTInt64         = "LISTInt64"
	LISTBytes         = "LISTBytes"
	STRUCT            = "STRUCT"
	STRUCT_END        = "STRUCT_END" //别用, 只会在struct解码的时候出现
	ZERO_TAG          = "ZERO_TAG"
	BOOL              = "BOOL"
	BYTES             = "BYTES"
)

type JceSection struct {
	JceType string
	Data    interface{}
}

func encodeHeadByte(jceId uint8, jceType uint8) []byte {
	if jceId < 15 {
		return []byte{jceId<<4 | jceType}
	} else {
		return []byte{0xF0 | jceType, jceId}
	}
}

func decodeHeadByte(firstHeaderByte byte) (uint8, uint8, bool) { //jceId jceType idInNexByte
	jceType := uint8(firstHeaderByte & 0x0F)
	jceId := uint8(firstHeaderByte >> 4)
	if jceId == 0x0F {
		return 0, jceType, true
	}
	return jceId, jceType, false
}
