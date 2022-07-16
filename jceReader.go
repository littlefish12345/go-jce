package gojce

import (
	"io"
)

type JceReader struct {
	data    []byte
	pointer uint64
}

func NewJceReader(data []byte) *JceReader {
	return &JceReader{data: data, pointer: 0}
}

func (jceReader *JceReader) Read(num uint64) ([]byte, error) {
	if jceReader.pointer == uint64(len(jceReader.data)) {
		return nil, io.EOF
	}
	if num+jceReader.pointer > uint64(len(jceReader.data)) {
		jceReader.pointer = uint64(len(jceReader.data))
		return jceReader.data[jceReader.pointer:], nil
	}
	jceReader.pointer += num
	return jceReader.data[jceReader.pointer-num : jceReader.pointer], nil
}

func (jceReader *JceReader) ReadHead() (uint8, uint8, error) { //jceId jceType
	if uint64(len(jceReader.data)) <= jceReader.pointer {
		return 0, 0, io.EOF
	}
	jceId, jceType, idInNexByte := JceSectionDecodeHeadByte(jceReader.data[jceReader.pointer])
	if idInNexByte {
		if uint64(len(jceReader.data)) <= jceReader.pointer+1 {
			return 0, 0, io.EOF
		}
		jceId = jceReader.data[jceReader.pointer+1]
	}
	return jceId, jceType, nil
}

func (jceReader *JceReader) SkipHead() {
	_, _, idInNexByte := JceSectionDecodeHeadByte(jceReader.data[jceReader.pointer])
	if idInNexByte {
		jceReader.pointer += 2
		return
	}
	jceReader.pointer += 1
}

func (jceReader *JceReader) SkipOneId() error {
	_, jceType, err := jceReader.ReadHead()
	if err != nil {
		return err
	}
	//fmt.Println(jceReader.pointer, jceId, jceType)
	jceReader.SkipHead()
	if jceType == 0 {
		jceReader.pointer += 1
	} else if jceType == 1 {
		jceReader.pointer += 2
	} else if jceType == 2 {
		jceReader.pointer += 4
	} else if jceType == 3 {
		jceReader.pointer += 8
	} else if jceType == 4 {
		jceReader.pointer += 4
	} else if jceType == 5 {
		jceReader.pointer += 8
	} else if jceType == 6 {
		lengthBytes, _ := jceReader.Read(1)
		jceReader.pointer += uint64(lengthBytes[0])
	} else if jceType == 7 {
		lengthBytes, _ := jceReader.Read(4)
		jceReader.pointer += uint64(BytesToInt32(lengthBytes))
	} else if jceType == 8 {
		length := JceSectionInt32FromBytes(jceReader) * 2
		for i := int32(0); i < length; i++ {
			jceReader.SkipOneId()
		}
	} else if jceType == 9 {
		length := JceSectionInt32FromBytes(jceReader)
		for i := int32(0); i < length; i++ {
			jceReader.SkipOneId()
		}
	} else if jceType == 10 {
		jceReader.SkipToStructEnd()
	} else if jceType == 13 {
		jceReader.SkipHead()
		jceReader.pointer += uint64(JceSectionInt32FromBytes(jceReader))
	}
	return nil
}

func (jceReader *JceReader) SkipToStructEnd() error {
	//fmt.Println("struct start")
	for {
		_, jceType, err := jceReader.ReadHead()
		if err != nil {
			return err
		}
		if jceType == 11 {
			jceReader.SkipHead()
			break
		}
		jceReader.SkipOneId()
	}
	//fmt.Println("struct end")
	return nil
}

func (jceReader *JceReader) SkipToId(jceTargetId uint8) error {
	for {
		jceId, _, err := jceReader.ReadHead()
		if err != nil {
			return err
		}
		//fmt.Println(jceId)
		if jceId >= jceTargetId {
			return nil
		}
		jceReader.SkipOneId()
	}
}

func (jceReader *JceReader) ReadJceStructByte() ([]byte, error) {
	_, jceType, err := jceReader.ReadHead()
	if err != nil {
		return nil, err
	}
	if jceType != 10 {
		return nil, nil
	}
	startPos := jceReader.pointer
	jceReader.SkipHead()
	err = jceReader.SkipToStructEnd()
	if err != nil {
		return nil, err
	}
	endPos := jceReader.pointer
	return jceReader.data[startPos:endPos], nil
}
