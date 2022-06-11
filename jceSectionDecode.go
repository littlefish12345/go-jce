package gojce

import (
	"bytes"
)

func jceSectionType0FromBytes(readBuffer *bytes.Buffer) (string, int8, error) { //jceType=0 int8
	buffer := make([]byte, 1)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", 0, err
	}
	return INT8, int8(buffer[0]), nil
}

func jceSectionType1FromBytes(readBuffer *bytes.Buffer) (string, int16, error) { //jceType=1 int16
	buffer := make([]byte, 2)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", 0, err
	}
	return INT16, BytesToInt16(buffer), nil
}

func jceSectionType2FromBytes(readBuffer *bytes.Buffer) (string, int32, error) { //jceType=2 int32
	buffer := make([]byte, 4)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", 0, err
	}
	return INT32, BytesToInt32(buffer), nil
}

func jceSectionType3FromBytes(readBuffer *bytes.Buffer) (string, int64, error) { //jceType=3 int64
	buffer := make([]byte, 8)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", 0, err
	}
	return INT64, BytesToInt64(buffer), nil
}

func jceSectionType4FromBytes(readBuffer *bytes.Buffer) (string, float32, error) { //jceType=4 float32
	buffer := make([]byte, 4)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", 0, err
	}
	return FLOAT32, BytesToFloat32(buffer), nil
}

func jceSectionType5FromBytes(readBuffer *bytes.Buffer) (string, float64, error) { //jceType=5 float64
	buffer := make([]byte, 8)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", 0, err
	}
	return FLOAT64, BytesToFloat64(buffer), nil
}

func jceSectionType6or7FromBytes(readBuffer *bytes.Buffer, jceType uint8) (string, string, error) { //jceType=6 or 7 string1/4
	var buffer []byte
	var length uint64
	if jceType == 6 {
		buffer = make([]byte, 1)
		_, err := readBuffer.Read(buffer)
		if err != nil {
			return "", "", err
		}
		length = uint64(buffer[0])
	} else {
		buffer = make([]byte, 4)
		_, err := readBuffer.Read(buffer)
		if err != nil {
			return "", "", err
		}
		length = uint64(BytesToInt64(buffer))
	}
	buffer = make([]byte, length)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", "", err
	}
	return STRING, string(buffer), nil
}

func jceSectionType8FromBytes(readBuffer *bytes.Buffer) (string, interface{}, error) { //jceType=8 map
	buffer := make([]byte, 1)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", nil, err
	}
	_, jceType, _ := decodeHeadByte(buffer[0])
	var jceSectionData interface{}
	var length int32
	if jceType == 12 {
		return MAP, nil, nil
	} else if jceType == 0 {
		_, jceSectionData, err = jceSectionType0FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		length = int32(jceSectionData.(int8))
	} else if jceType == 1 {
		_, jceSectionData, err = jceSectionType1FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		length = int32(jceSectionData.(int16))
	} else if jceType == 2 {
		_, jceSectionData, err = jceSectionType2FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		length = jceSectionData.(int32)
	} else {
		return "", nil, ErrorUnknowJceTypeErr
	}

	var returnMap interface{}
	var innerMapType string
	var value interface{}
	mapType := ""
	_, err = readBuffer.Read(buffer)
	if err != nil {
		return "", nil, err
	}
	_, jceType, _ = decodeHeadByte(buffer[0])
	if jceType != 6 && jceType != 7 {
		return "", nil, ErrorUnknowJceTypeErr
	}
	_, key, err := jceSectionType6or7FromBytes(readBuffer, jceType)
	if err != nil {
		return "", nil, err
	}

	_, err = readBuffer.Read(buffer)
	if err != nil {
		return "", nil, err
	}
	_, jceType, _ = decodeHeadByte(buffer[0])
	if jceType == 6 || jceType == 7 {
		returnMap = make(map[string]string)
		_, value, err = jceSectionType6or7FromBytes(readBuffer, jceType)
		if err != nil {
			return "", nil, err
		}
		returnMap.(map[string]interface{})[key] = value
		mapType = MAPStrStr
	} else if jceType == 12 {
		returnMap = make(map[string][]byte)
		_, value, err = jceSectionType13FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		returnMap.(map[string]interface{})[key] = value
		mapType = MAPStrBytes
	} else if jceType == 8 {
		returnMap = make(map[string]map[string][]byte)
		innerMapType, value, err = jceSectionType8FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		if innerMapType != MAPStrBytes {
			return "", nil, ErrorUnknowJceTypeErr
		}
		returnMap.(map[string]interface{})[key] = value
		mapType = MAPStrMAPStrBytes
	}

	for i := int32(1); i < length; i++ {
		_, err = readBuffer.Read(buffer)
		if err != nil {
			return "", nil, err
		}
		_, jceType, _ = decodeHeadByte(buffer[0])
		if jceType != 6 && jceType != 7 {
			return "", nil, ErrorUnknowJceTypeErr
		}
		_, key, err = jceSectionType6or7FromBytes(readBuffer, jceType)
		if err != nil {
			return "", nil, err
		}

		_, err = readBuffer.Read(buffer)
		if err != nil {
			return "", nil, err
		}
		_, jceType, _ = decodeHeadByte(buffer[0])
		if mapType == MAPStrStr {
			if jceType != 6 && jceType != 7 {
				return "", nil, ErrorUnknowJceTypeErr
			}
			_, value, err = jceSectionType6or7FromBytes(readBuffer, jceType)
			if err != nil {
				return "", nil, err
			}
		} else if mapType == MAPStrBytes {
			if jceType != 12 {
				return "", nil, ErrorUnknowJceTypeErr
			}
			_, value, err = jceSectionType13FromBytes(readBuffer)
			if err != nil {
				return "", nil, err
			}
		} else {
			if jceType != 8 {
				return "", nil, ErrorUnknowJceTypeErr
			}
			innerMapType, value, err = jceSectionType8FromBytes(readBuffer)
			if err != nil {
				return "", nil, err
			}
			if innerMapType != MAPStrBytes {
				return "", nil, ErrorUnknowJceTypeErr
			}
		}
		returnMap.(map[string]interface{})[key] = value
	}
	return mapType, returnMap, nil
}

func jceSectionType9FromBytes(readBuffer *bytes.Buffer) (string, interface{}, error) { //jceType=9 list
	buffer := make([]byte, 1)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", nil, err
	}
	_, jceType, _ := decodeHeadByte(buffer[0])
	var jceSectionData interface{}
	var length int32
	if jceType == 12 {
		return LIST, nil, nil
	} else if jceType == 0 {
		_, jceSectionData, err = jceSectionType0FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		length = int32(jceSectionData.(int8))
	} else if jceType == 1 {
		_, jceSectionData, err = jceSectionType1FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		length = int32(jceSectionData.(int16))
	} else if jceType == 2 {
		_, jceSectionData, err = jceSectionType2FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		length = jceSectionData.(int32)
	} else {
		return "", nil, ErrorUnknowJceTypeErr
	}

	var returnList []interface{}
	var value interface{}
	listType := ""
	_, err = readBuffer.Read(buffer)
	if err != nil {
		return "", nil, err
	}
	_, jceType, _ = decodeHeadByte(buffer[0])
	if jceType == 0 || jceType == 1 || jceType == 2 || jceType == 3 || jceType == 12 {
		_, value, err = jceSectionType6or7FromBytes(readBuffer, jceType)
		if err != nil {
			return "", nil, err
		}
		if jceType == 0 {
			_, jceSectionData, err = jceSectionType0FromBytes(readBuffer)
			if err != nil {
				return "", nil, err
			}
			value = int64(jceSectionData.(int8))
		} else if jceType == 1 {
			_, jceSectionData, err = jceSectionType1FromBytes(readBuffer)
			if err != nil {
				return "", nil, err
			}
			value = int64(jceSectionData.(int16))
		} else if jceType == 2 {
			_, jceSectionData, err = jceSectionType2FromBytes(readBuffer)
			if err != nil {
				return "", nil, err
			}
			value = int64(jceSectionData.(int32))
		} else if jceType == 3 {
			_, jceSectionData, err = jceSectionType3FromBytes(readBuffer)
			if err != nil {
				return "", nil, err
			}
			value = int64(jceSectionData.(int64))
		} else {
			value = int64(0)
		}
		returnList = append(returnList, value)
		listType = LISTInt64
	} else if jceType == 13 {
		_, value, err = jceSectionType13FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		returnList = append(returnList, value)
		listType = LISTBytes
	} else {
		return "", nil, ErrorUnknowJceTypeErr
	}

	for i := int32(1); i < length; i++ {
		_, err = readBuffer.Read(buffer)
		if err != nil {
			return "", nil, err
		}
		_, jceType, _ = decodeHeadByte(buffer[0])
		if listType == LISTInt64 {
			if jceType != 0 && jceType != 1 && jceType != 2 && jceType != 3 && jceType != 12 {
				return "", nil, ErrorUnknowJceTypeErr
			}
			if jceType == 0 {
				_, jceSectionData, err = jceSectionType0FromBytes(readBuffer)
				if err != nil {
					return "", nil, err
				}
				value = int64(jceSectionData.(int8))
			} else if jceType == 1 {
				_, jceSectionData, err = jceSectionType1FromBytes(readBuffer)
				if err != nil {
					return "", nil, err
				}
				value = int64(jceSectionData.(int16))
			} else if jceType == 2 {
				_, jceSectionData, err = jceSectionType2FromBytes(readBuffer)
				if err != nil {
					return "", nil, err
				}
				value = int64(jceSectionData.(int32))
			} else if jceType == 3 {
				_, jceSectionData, err = jceSectionType3FromBytes(readBuffer)
				if err != nil {
					return "", nil, err
				}
				value = int64(jceSectionData.(int64))
			} else {
				value = int64(0)
			}
		} else {
			if jceType != 13 {
				return "", nil, ErrorUnknowJceTypeErr
			}
			_, value, err = jceSectionType13FromBytes(readBuffer)
			if err != nil {
				return "", nil, err
			}
		}
		returnList = append(returnList, value)
	}
	return listType, returnList, nil
}

func jceSectionType13FromBytes(readBuffer *bytes.Buffer) (string, []byte, error) { //jceType=13
	buffer := make([]byte, 1)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return "", nil, err
	}
	buffer = make([]byte, 1)
	_, err = readBuffer.Read(buffer)
	if err != nil {
		return "", nil, err
	}
	_, jceType, _ := decodeHeadByte(buffer[0])
	var jceSectionData interface{}
	var length uint32
	if jceType == 12 {
		return BYTES, []byte{}, nil
	} else if jceType == 0 {
		_, jceSectionData, err = jceSectionType0FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		length = uint32(jceSectionData.(int8))
	} else if jceType == 1 {
		_, jceSectionData, err = jceSectionType1FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		length = uint32(jceSectionData.(int16))
	} else if jceType == 2 {
		_, jceSectionData, err = jceSectionType2FromBytes(readBuffer)
		if err != nil {
			return "", nil, err
		}
		length = uint32(jceSectionData.(int32))
	}
	buffer = make([]byte, length)
	_, err = readBuffer.Read(buffer)
	if err != nil {
		return "", nil, err
	}
	return BYTES, buffer, nil
}

func (jceSection *JceSection) Decode(readBuffer *bytes.Buffer) (uint8, error) {
	buffer := make([]byte, 1)
	_, err := readBuffer.Read(buffer)
	if err != nil {
		return 0, err
	}
	jceId, jceType, idInNexByte := decodeHeadByte(buffer[0])
	if idInNexByte {
		_, err := readBuffer.Read(buffer)
		if err != nil {
			return 0, err
		}
		jceId = buffer[0]
	}
	var jceSectionType string
	var jceSectionData interface{}
	if jceType == 0 {
		jceSectionType, jceSectionData, err = jceSectionType0FromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
	} else if jceType == 1 {
		jceSectionType, jceSectionData, err = jceSectionType1FromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
	} else if jceType == 2 {
		jceSectionType, jceSectionData, err = jceSectionType2FromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
	} else if jceType == 3 {
		jceSectionType, jceSectionData, err = jceSectionType3FromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
	} else if jceType == 4 {
		jceSectionType, jceSectionData, err = jceSectionType4FromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
	} else if jceType == 5 {
		jceSectionType, jceSectionData, err = jceSectionType5FromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
	} else if jceType == 6 || jceType == 7 {
		jceSectionType, jceSectionData, err = jceSectionType6or7FromBytes(readBuffer, jceType)
		if err != nil {
			return 0, err
		}
	} else if jceType == 8 {
		jceSectionType, jceSectionData, err = jceSectionType8FromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
	} else if jceType == 9 {
		jceSectionType, jceSectionData, err = jceSectionType9FromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
	} else if jceType == 10 {
		var data map[uint8]*JceSection
		jceSectionType, data, err = jceStructFromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
		jceSectionData = newJceStruct()
		jceSectionData.(*JceStruct).structMap = data
	} else if jceType == 11 {
		jceSectionType = STRUCT_END
	} else if jceType == 12 {
		jceSectionType = ZERO_TAG
	} else if jceType == 13 {
		jceSectionType, jceSectionData, err = jceSectionType13FromBytes(readBuffer)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, ErrorUnknowJceTypeErr
	}
	*jceSection = JceSection{JceType: jceSectionType, Data: jceSectionData}
	return jceId, nil
}
