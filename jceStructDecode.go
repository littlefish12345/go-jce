package gojce

import (
	"io"
)

func jceStructDecode(readBuffer *JceReader) (map[uint8]*JceSection, error) {
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
		if jceSection.JceType != STRUCT && jceSection.JceType != STRUCT_END {
			if _, ok := jceStruct[jceId]; !ok {
				jceStruct[jceId] = jceSection
			}
		}
	}
	return jceStruct, nil
}

func jceStructFromBytes(readBuffer *JceReader) map[uint8]*JceSection {
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 10 {
		return nil
	}
	readBuffer.SkipHead()
	jceStruct := make(map[uint8]*JceSection)
	var jceSection *JceSection
	for i := uint8(0); ; i++ {
		jceSection = &JceSection{}
		jceId, err := jceSection.Decode(readBuffer)
		if err != nil {
			return nil
		}
		if jceSection.JceType == STRUCT_END {
			break
		}
		jceStruct[jceId] = jceSection
	}
	return jceStruct
}

func (jceStruct *JceStruct) Decode(data []byte) error {
	readBuffer := NewJceReader(data)
	var err error
	jceStruct.structMap, err = jceStructDecode(readBuffer)
	return err
}

func (jceStruct *JceStruct) DecodeWithLength(data []byte) error {
	length := BytesToInt32(data[0:4])
	data = data[4:length]
	readBuffer := NewJceReader(data)
	var err error
	jceStruct.structMap, err = jceStructDecode(readBuffer)
	return err
}

func (jceStruct *JceStruct) FromBytes(readBuffer *JceReader) error {
	/*
		buffer := make([]byte, 1)
		_, err = readBuffer.Read(buffer)
		if err != nil {
			return 0, err
		}
		jceId, jceType, idInNexByte := JceSectionDecodeHeadByte(buffer[0])
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
	*/
	var err error
	jceStruct.structMap = jceStructFromBytes(readBuffer)
	return err
}
