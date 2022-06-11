package gojce

import (
	"bytes"
)

func jceStructEncode(jceStruct map[uint8]*JceSection) ([]byte, error) {
	buffer := new(bytes.Buffer)
	for i := uint8(0); ; i++ {
		if jceSection, ok := jceStruct[i]; ok {
			data, err := jceSection.Encode(i)
			if err != nil {
				return nil, err
			}
			buffer.Write(data)
		} else {
			if i != 0 {
				break
			}
		}
	}
	return buffer.Bytes(), nil
}

func jceStructToBytes(jceId uint8, jceStruct map[uint8]*JceSection) ([]byte, error) {
	buffer := new(bytes.Buffer)
	buffer.Write(jceSectionStructStartToBytes(jceId))
	data, err := jceStructEncode(jceStruct)
	if err != nil {
		return nil, err
	}
	buffer.Write(data)
	buffer.Write(jceSectionStructEndToBytes(jceId))
	return buffer.Bytes(), nil
}

func (jceStruct *JceStruct) Encode() ([]byte, error) {
	return jceStructEncode(jceStruct.structMap)
}

func (jceStruct *JceStruct) EncodeWithLength() ([]byte, error) {
	data, err := jceStruct.Encode()
	data = append(Int32ToBytes(int32(len(data)+4)), data...)
	return data, err
}

func (jceStruct *JceStruct) ToBytes(jceId uint8) ([]byte, error) {
	return jceStructToBytes(jceId, jceStruct.structMap)
}
