package gojce

import (
	"bytes"
	"fmt"
)

func jceSectionInt8ToBytes(jceId uint8, data int8) []byte { //jceType=0
	if data == 0x00 {
		return jceSectionZeroTagToBytes(jceId)
	} else {
		return append(encodeHeadByte(jceId, 0), uint8(data))
	}
}

func jceSectionByteToBytes(jceId uint8, data uint8) []byte { //jceType=0
	return jceSectionInt8ToBytes(jceId, int8(data))
}

func jceSectionInt16ToBytes(jceId uint8, data int16) []byte { //jceType=1
	return append(encodeHeadByte(jceId, 1), Int16ToBytes(data)...)
}

func jceSectionInt32ToBytes(jceId uint8, data int32) []byte { //jceType=2
	return append(encodeHeadByte(jceId, 2), Int32ToBytes(data)...)
}

func jceSectionInt64ToBytes(jceId uint8, data int64) []byte { //jceType=3
	return append(encodeHeadByte(jceId, 3), Int64ToBytes(data)...)
}

func jceSectionIntToBytes(jceId uint8, data int64) []byte { //jceType=0 or 1 or 2 or 3
	if -128 <= data && data <= 127 {
		return jceSectionInt8ToBytes(jceId, int8(data))
	} else if -32768 <= data && data <= 32767 {
		return jceSectionInt16ToBytes(jceId, int16(data))
	} else if -2147483648 <= data && data <= 2147483647 {
		return jceSectionInt32ToBytes(jceId, int32(data))
	}
	return jceSectionInt64ToBytes(jceId, data)
}

func jceSectionFloat32ToBytes(jceId uint8, data float32) []byte { //jceType=4
	return append(encodeHeadByte(jceId, 4), Float32ToByte(data)...)
}

func jceSectionFloat64ToBytes(jceId uint8, data float64) []byte { //jceType=5
	return append(encodeHeadByte(jceId, 5), Float64ToByte(data)...)
}

func jceSectionString1ToBytes(jceId uint8, data []byte) []byte { //jceType=6
	return append(append(encodeHeadByte(jceId, 6), byte(len(data))), data...)
}

func jceSectionString4ToBytes(jceId uint8, data []byte) []byte { //jceType=7
	return append(append(encodeHeadByte(jceId, 7), Int32ToBytes(int32(len(data)))...), data...)
}

func jceSectionStringToBytes(jceId uint8, data string) []byte { //jceType=6 or 7
	dataByte := []byte(data)
	if len(data) < 256 {
		return jceSectionString1ToBytes(jceId, dataByte)
	}
	return jceSectionString4ToBytes(jceId, dataByte)
}

func JceSectionMapStrStrToBytes(jceId uint8, data map[string]string) []byte { //jceType=8
	buffer := new(bytes.Buffer)
	buffer.Write(encodeHeadByte(jceId, 8))
	if data == nil {
		buffer.Write(jceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(jceSectionInt32ToBytes(0, int32(len(data))))
	for key, value := range data {
		buffer.Write(jceSectionStringToBytes(0, key))
		buffer.Write(jceSectionStringToBytes(1, value))
	}
	return buffer.Bytes()
}

func JceSectionMapStrBytesToBytes(jceId uint8, data map[string][]byte) []byte { //jceType=8
	buffer := new(bytes.Buffer)
	buffer.Write(encodeHeadByte(jceId, 8))
	if data == nil {
		buffer.Write(jceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(jceSectionInt32ToBytes(0, int32(len(data))))
	for key, value := range data {
		buffer.Write(jceSectionStringToBytes(0, key))
		buffer.Write(jceSectionBytesToBytes(1, value))
	}
	return buffer.Bytes()
}

func JceSectionMapStrMapStrBytesToBytes(jceId uint8, data map[string]map[string][]byte) []byte { //jceType=8
	buffer := new(bytes.Buffer)
	buffer.Write(encodeHeadByte(jceId, 8))
	if data == nil {
		buffer.Write(jceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(jceSectionInt32ToBytes(0, int32(len(data))))
	for key, value := range data {
		buffer.Write(jceSectionStringToBytes(0, key))
		buffer.Write(JceSectionMapStrBytesToBytes(1, value))
	}
	return buffer.Bytes()
}

func jceSectionListInt64ToBytes(jceId uint8, data []int64) []byte { //jceType=9
	buffer := new(bytes.Buffer)
	buffer.Write(encodeHeadByte(jceId, 9))
	if len(data) == 0 {
		buffer.Write(jceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(jceSectionInt32ToBytes(0, int32(len(data))))
	for _, value := range data {
		buffer.Write(jceSectionInt64ToBytes(1, value))
	}
	return buffer.Bytes()
}

func jceSectionListBytesToBytes(jceId uint8, data [][]byte) []byte { //jceType=9
	buffer := new(bytes.Buffer)
	buffer.Write(encodeHeadByte(jceId, 9))
	if len(data) == 0 {
		buffer.Write(jceSectionZeroTagToBytes(0))
		return buffer.Bytes()
	}
	buffer.Write(jceSectionInt32ToBytes(0, int32(len(data))))
	for _, value := range data {
		buffer.Write(jceSectionBytesToBytes(1, value))
	}
	return buffer.Bytes()
}

func jceSectionStructStartToBytes(jceId uint8) []byte { //jceType=10
	return encodeHeadByte(jceId, 10)
}

func jceSectionStructEndToBytes(jceId uint8) []byte { //jceType=11
	return encodeHeadByte(jceId, 11)
}

func jceSectionZeroTagToBytes(jceId uint8) []byte { //jceType=12
	return encodeHeadByte(jceId, 12)
}

func jceSectionBoolToBytes(jceId uint8, data bool) []byte { //jceType=0 or 12
	if data {
		return append(encodeHeadByte(jceId, 1), 0x01)
	} else {
		return jceSectionZeroTagToBytes(jceId)
	}
}

func jceSectionBytesToBytes(jceId uint8, data []byte) []byte { //jceType=13
	return append(encodeHeadByte(jceId, 13), append(encodeHeadByte(0, 0), append(jceSectionInt32ToBytes(0, int32(len(data))), data...)...)...)
}

func (jceSection JceSection) Encode(jceId uint8) ([]byte, error) {
	if jceSection.JceType != "" {
		if jceSection.JceType == INT8 {
			return jceSectionIntToBytes(jceId, int64(jceSection.Data.(int8))), nil
		} else if jceSection.JceType == BYTE {
			return jceSectionByteToBytes(jceId, jceSection.Data.(byte)), nil
		} else if jceSection.JceType == INT16 {
			return jceSectionIntToBytes(jceId, int64(jceSection.Data.(int16))), nil
		} else if jceSection.JceType == INT32 {
			return jceSectionIntToBytes(jceId, int64(jceSection.Data.(int32))), nil
		} else if jceSection.JceType == INT64 {
			return jceSectionIntToBytes(jceId, jceSection.Data.(int64)), nil
		} else if jceSection.JceType == FLOAT32 {
			return jceSectionFloat32ToBytes(jceId, jceSection.Data.(float32)), nil
		} else if jceSection.JceType == FLOAT64 {
			return jceSectionFloat64ToBytes(jceId, jceSection.Data.(float64)), nil
		} else if jceSection.JceType == STRING {
			return jceSectionStringToBytes(jceId, jceSection.Data.(string)), nil
		} else if jceSection.JceType == MAPStrStr {
			return JceSectionMapStrStrToBytes(jceId, jceSection.Data.(map[string]string)), nil
		} else if jceSection.JceType == MAPStrBytes {
			return JceSectionMapStrBytesToBytes(jceId, jceSection.Data.(map[string][]byte)), nil
		} else if jceSection.JceType == MAPStrMAPStrBytes {
			return JceSectionMapStrMapStrBytesToBytes(jceId, jceSection.Data.(map[string]map[string][]byte)), nil
		} else if jceSection.JceType == LISTInt64 {
			return jceSectionListInt64ToBytes(jceId, jceSection.Data.([]int64)), nil
		} else if jceSection.JceType == LISTBytes {
			return jceSectionListBytesToBytes(jceId, jceSection.Data.([][]byte)), nil
		} else if jceSection.JceType == BOOL {
			return jceSectionBoolToBytes(jceId, jceSection.Data.(bool)), nil
		} else if jceSection.JceType == BYTES {
			return jceSectionBytesToBytes(jceId, jceSection.Data.([]byte)), nil
		}
	}
	fmt.Println(jceSection.JceType)
	return nil, ErrorUnknowJceTypeErr
}
