package gojce

import (
	"bytes"
	"io"
)

func jceStructDecode(readBuffer *bytes.Buffer) (map[uint8]*JceSection, error) {
	jceStruct := make(map[uint8]*JceSection)
	var jceSection *JceSection
	for i := uint8(0); ; i++ {
		jceSection = &JceSection{}
		jceId, err := jceSection.Decode(readBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		jceStruct[jceId] = jceSection
	}
	return jceStruct, nil
}

func jceStructFromBytes(readBuffer *bytes.Buffer) (string, map[uint8]*JceSection, error) {
	jceStruct := make(map[uint8]*JceSection)
	var jceSection *JceSection
	for i := uint8(0); ; i++ {
		jceSection = &JceSection{}
		jceId, err := jceSection.Decode(readBuffer)
		if err != nil {
			return "", nil, err
		}
		if jceSection.JceType == STRUCT_END {
			break
		}
		jceStruct[jceId] = jceSection
	}
	return STRUCT, jceStruct, nil
}

func (jceStruct *JceStruct) Decode(data []byte) error {
	readBuffer := bytes.NewBuffer(data)
	var err error
	jceStruct.structMap, err = jceStructDecode(readBuffer)
	return err
}

func (jceStruct *JceStruct) DecodeWithLength(data []byte) error {
	length := BytesToInt32(data[0:4])
	data = data[4:length]
	readBuffer := bytes.NewBuffer(data)
	var err error
	jceStruct.structMap, err = jceStructDecode(readBuffer)
	return err
}

func (jceStruct *JceStruct) FromBytes(readBuffer *bytes.Buffer) (uint8, error) {
	var err error
	buffer := make([]byte, 1)
	_, err = readBuffer.Read(buffer)
	if err != nil {
		return 0, err
	}
	jceId, jceType, idInNexByte := decodeHeadByte(buffer[0])
	if jceType != 10 {
		return 0, ErrorUnknowJceTypeErr
	}
	if idInNexByte {
		_, err := readBuffer.Read(buffer)
		if err != nil {
			return 0, err
		}
		jceId = buffer[0]
	}
	_, jceStruct.structMap, err = jceStructFromBytes(readBuffer)
	return jceId, err
}
