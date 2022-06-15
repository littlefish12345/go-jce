package gojce

import (
	"bytes"
	"fmt"
)

func JceSectionInt8ToBytes(jceId uint8, data int8) []byte { //jceType=0
	if data == 0x00 {
		return JceSectionZeroTagToBytes(jceId)
	} else {
		return append(JceSectionEncodeHeadByte(jceId, 0), uint8(data))
	}
}

func JceSectionByteToBytes(jceId uint8, data uint8) []byte { //jceType=0
	return JceSectionInt8ToBytes(jceId, int8(data))
}

func JceSectionInt16ToBytes(jceId uint8, data int16) []byte { //jceType=1
	if -128 <= data && data <= 127 {
		return JceSectionInt8ToBytes(jceId, int8(data))
	}
	return append(JceSectionEncodeHeadByte(jceId, 1), Int16ToBytes(data)...)
}

func JceSectionInt32ToBytes(jceId uint8, data int32) []byte { //jceType=2
	if -128 <= data && data <= 127 {
		return JceSectionInt8ToBytes(jceId, int8(data))
	} else if -32768 <= data && data <= 32767 {
		return JceSectionInt16ToBytes(jceId, int16(data))
	}
	return append(JceSectionEncodeHeadByte(jceId, 2), Int32ToBytes(data)...)
}

func JceSectionInt64ToBytes(jceId uint8, data int64) []byte { //jceType=3
	if -128 <= data && data <= 127 {
		return JceSectionInt8ToBytes(jceId, int8(data))
	} else if -32768 <= data && data <= 32767 {
		return JceSectionInt16ToBytes(jceId, int16(data))
	} else if -2147483648 <= data && data <= 2147483647 {
		return JceSectionInt32ToBytes(jceId, int32(data))
	}
	return append(JceSectionEncodeHeadByte(jceId, 3), Int64ToBytes(data)...)
}

func JceSectionIntToBytes(jceId uint8, data int64) []byte { //jceType=0 or 1 or 2 or 3
	if -128 <= data && data <= 127 {
		return JceSectionInt8ToBytes(jceId, int8(data))
	} else if -32768 <= data && data <= 32767 {
		return JceSectionInt16ToBytes(jceId, int16(data))
	} else if -2147483648 <= data && data <= 2147483647 {
		return JceSectionInt32ToBytes(jceId, int32(data))
	}
	return JceSectionInt64ToBytes(jceId, data)
}

func JceSectionFloat32ToBytes(jceId uint8, data float32) []byte { //jceType=4
	return append(JceSectionEncodeHeadByte(jceId, 4), Float32ToByte(data)...)
}

func JceSectionFloat64ToBytes(jceId uint8, data float64) []byte { //jceType=5
	return append(JceSectionEncodeHeadByte(jceId, 5), Float64ToByte(data)...)
}

func JceSectionString1ToBytes(jceId uint8, data []byte) []byte { //jceType=6
	return append(append(JceSectionEncodeHeadByte(jceId, 6), byte(len(data))), data...)
}

func JceSectionString4ToBytes(jceId uint8, data []byte) []byte { //jceType=7
	return append(append(JceSectionEncodeHeadByte(jceId, 7), Int32ToBytes(int32(len(data)))...), data...)
}

func JceSectionStringToBytes(jceId uint8, data string) []byte { //jceType=6 or 7
	dataByte := []byte(data)
	if len(data) < 256 {
		return JceSectionString1ToBytes(jceId, dataByte)
	}
	return JceSectionString4ToBytes(jceId, dataByte)
}

func JceSectionMapStrStrToBytes(jceId uint8, data map[string]string) []byte { //jceType=8
	buffer := new(bytes.Buffer)
	buffer.Write(JceSectionEncodeHeadByte(jceId, 8))
	if data == nil {
		buffer.Write(JceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(JceSectionInt32ToBytes(0, int32(len(data))))
	for key, value := range data {
		buffer.Write(JceSectionStringToBytes(0, key))
		buffer.Write(JceSectionStringToBytes(1, value))
	}
	return buffer.Bytes()
}

func JceSectionMapStrBytesToBytes(jceId uint8, data map[string][]byte) []byte { //jceType=8
	buffer := new(bytes.Buffer)
	buffer.Write(JceSectionEncodeHeadByte(jceId, 8))
	if data == nil {
		buffer.Write(JceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(JceSectionInt32ToBytes(0, int32(len(data))))
	for key, value := range data {
		buffer.Write(JceSectionStringToBytes(0, key))
		buffer.Write(JceSectionBytesToBytes(1, value))
	}
	return buffer.Bytes()
}

func JceSectionMapStrMapStrBytesToBytes(jceId uint8, data map[string]map[string][]byte) []byte { //jceType=8
	buffer := new(bytes.Buffer)
	buffer.Write(JceSectionEncodeHeadByte(jceId, 8))
	if data == nil {
		buffer.Write(JceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(JceSectionInt32ToBytes(0, int32(len(data))))
	for key, value := range data {
		buffer.Write(JceSectionStringToBytes(0, key))
		buffer.Write(JceSectionMapStrBytesToBytes(1, value))
	}
	return buffer.Bytes()
}

func JceSectionListInt64ToBytes(jceId uint8, data []int64) []byte { //jceType=9
	buffer := new(bytes.Buffer)
	buffer.Write(JceSectionEncodeHeadByte(jceId, 9))
	if len(data) == 0 {
		buffer.Write(JceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(JceSectionInt32ToBytes(0, int32(len(data))))
	for _, value := range data {
		buffer.Write(JceSectionInt64ToBytes(1, value))
	}
	return buffer.Bytes()
}

func JceSectionListBytesToBytes(jceId uint8, data [][]byte) []byte { //jceType=9
	buffer := new(bytes.Buffer)
	buffer.Write(JceSectionEncodeHeadByte(jceId, 9))
	if len(data) == 0 {
		buffer.Write(JceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(JceSectionInt32ToBytes(0, int32(len(data))))
	for _, value := range data {
		buffer.Write(JceSectionBytesToBytes(1, value))
	}
	return buffer.Bytes()
}

func JceSectionStructStartToBytes(jceId uint8) []byte { //jceType=10
	return JceSectionEncodeHeadByte(jceId, 10)
}

func JceSectionStructEndToBytes(jceId uint8) []byte { //jceType=11
	return JceSectionEncodeHeadByte(jceId, 11)
}

func JceSectionZeroTagToBytes(jceId uint8) []byte { //jceType=12
	return JceSectionEncodeHeadByte(jceId, 12)
}

func JceSectionBoolToBytes(jceId uint8, data bool) []byte { //jceType=0 or 12
	if data {
		return append(JceSectionEncodeHeadByte(jceId, 1), 0x01)
	} else {
		return JceSectionZeroTagToBytes(jceId)
	}
}

func JceSectionBytesToBytes(jceId uint8, data []byte) []byte { //jceType=13
	return append(JceSectionEncodeHeadByte(jceId, 13), append(JceSectionEncodeHeadByte(0, 0), append(JceSectionInt32ToBytes(0, int32(len(data))), data...)...)...)
}

func (jceSection JceSection) Encode(jceId uint8) ([]byte, error) {
	if jceSection.JceType != "" {
		if jceSection.JceType == INT8 {
			return JceSectionIntToBytes(jceId, int64(jceSection.Data.(int8))), nil
		} else if jceSection.JceType == BYTE {
			return JceSectionByteToBytes(jceId, jceSection.Data.(byte)), nil
		} else if jceSection.JceType == INT16 {
			return JceSectionIntToBytes(jceId, int64(jceSection.Data.(int16))), nil
		} else if jceSection.JceType == INT32 {
			return JceSectionIntToBytes(jceId, int64(jceSection.Data.(int32))), nil
		} else if jceSection.JceType == INT64 {
			return JceSectionIntToBytes(jceId, jceSection.Data.(int64)), nil
		} else if jceSection.JceType == FLOAT32 {
			return JceSectionFloat32ToBytes(jceId, jceSection.Data.(float32)), nil
		} else if jceSection.JceType == FLOAT64 {
			return JceSectionFloat64ToBytes(jceId, jceSection.Data.(float64)), nil
		} else if jceSection.JceType == STRING {
			return JceSectionStringToBytes(jceId, jceSection.Data.(string)), nil
		} else if jceSection.JceType == MAPStrStr {
			return JceSectionMapStrStrToBytes(jceId, jceSection.Data.(map[string]string)), nil
		} else if jceSection.JceType == MAPStrBytes {
			return JceSectionMapStrBytesToBytes(jceId, jceSection.Data.(map[string][]byte)), nil
		} else if jceSection.JceType == MAPStrMAPStrBytes {
			return JceSectionMapStrMapStrBytesToBytes(jceId, jceSection.Data.(map[string]map[string][]byte)), nil
		} else if jceSection.JceType == LISTInt64 {
			return JceSectionListInt64ToBytes(jceId, jceSection.Data.([]int64)), nil
		} else if jceSection.JceType == LISTBytes {
			return JceSectionListBytesToBytes(jceId, jceSection.Data.([][]byte)), nil
		} else if jceSection.JceType == BOOL {
			return JceSectionBoolToBytes(jceId, jceSection.Data.(bool)), nil
		} else if jceSection.JceType == BYTES {
			return JceSectionBytesToBytes(jceId, jceSection.Data.([]byte)), nil
		} else if jceSection.JceType == STRUCT {
			return jceSection.Data.(*JceStruct).ToBytes(jceId)
		} else {
			return nil, ErrorUnknowJceTypeErr
		}
	}
	fmt.Println(jceSection.JceType)
	return nil, ErrorUnknowJceTypeErr
}
