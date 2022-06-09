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
	MAPStrStr         = "MAPStrStr"
	MAPStrBytes       = "MAPStrBytes"
	MAPStrMAPStrBytes = "MAPStrMAPStrBytes"
	LISTInt64         = "LISTInt64"
	LISTBytes         = "LISTBytes"
	STRUCT_START      = "STRUCT_START"
	STRUCT_END        = "STRUCT_END"
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
