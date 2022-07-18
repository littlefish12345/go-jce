package gojce

import "fmt"

func JceSectionType0FromBytes(readBuffer *JceReader) int8 { //jceType=0 int8
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 0 {
		return 0
	}
	readBuffer.SkipHead()
	buffer, _ := readBuffer.Read(1)
	if len(buffer) != 1 {
		return 0
	}
	return int8(buffer[0])
}

func JceSectionType1FromBytes(readBuffer *JceReader) int16 { //jceType=1 int16
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 1 {
		return 0
	}
	readBuffer.SkipHead()
	buffer, _ := readBuffer.Read(2)
	if len(buffer) != 2 {
		return 0
	}
	return BytesToInt16(buffer)
}

func JceSectionType2FromBytes(readBuffer *JceReader) int32 { //jceType=2 int32
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 2 {
		return 0
	}
	readBuffer.SkipHead()
	buffer, _ := readBuffer.Read(4)
	if len(buffer) != 4 {
		return 0
	}
	return BytesToInt32(buffer)
}

func JceSectionType3FromBytes(readBuffer *JceReader) int64 { //jceType=3 int64
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 3 {
		return 0
	}
	readBuffer.SkipHead()
	buffer, _ := readBuffer.Read(8)
	if len(buffer) != 8 {
		return 0
	}
	return BytesToInt64(buffer)
}

func JceSectionInt8FromBytes(readBuffer *JceReader) int8 {
	_, jceType, _ := readBuffer.ReadHead()
	if jceType == 12 {
		return 0
	} else if jceType == 0 {
		return JceSectionType0FromBytes(readBuffer)
	}
	return 0
}

func JceSectionByteFromBytes(readBuffer *JceReader) byte {
	_, jceType, _ := readBuffer.ReadHead()
	if jceType == 12 {
		return 0
	} else if jceType == 0 {
		return byte(JceSectionType0FromBytes(readBuffer))
	}
	return 0
}

func JceSectionInt16FromBytes(readBuffer *JceReader) int16 {
	_, jceType, _ := readBuffer.ReadHead()
	if jceType == 0 {
		return int16(JceSectionType0FromBytes(readBuffer))
	} else if jceType == 1 {
		return JceSectionType1FromBytes(readBuffer)
	} else if jceType == 12 {
		return 0
	}
	return 0
}

func JceSectionInt32FromBytes(readBuffer *JceReader) int32 {
	_, jceType, _ := readBuffer.ReadHead()
	if jceType == 0 {
		return int32(JceSectionType0FromBytes(readBuffer))
	} else if jceType == 1 {
		return int32(JceSectionType1FromBytes(readBuffer))
	} else if jceType == 2 {
		return JceSectionType2FromBytes(readBuffer)
	} else if jceType == 12 {
		return 0
	}
	return 0
}

func JceSectionInt64FromBytes(readBuffer *JceReader) int64 {
	_, jceType, _ := readBuffer.ReadHead()
	if jceType == 0 {
		return int64(JceSectionType0FromBytes(readBuffer))
	} else if jceType == 1 {
		return int64(JceSectionType1FromBytes(readBuffer))
	} else if jceType == 2 {
		return int64(JceSectionType2FromBytes(readBuffer))
	} else if jceType == 3 {
		return JceSectionType3FromBytes(readBuffer)
	} else if jceType == 12 {
		return 0
	}
	return 0
}

func JceSectionType4FromBytes(readBuffer *JceReader) float32 { //jceType=4 float32
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 4 {
		return 0
	}
	readBuffer.SkipHead()
	buffer, _ := readBuffer.Read(4)
	if len(buffer) != 4 {
		return 0
	}
	return BytesToFloat32(buffer)
}

func JceSectionType5FromBytes(readBuffer *JceReader) float64 { //jceType=5 float64
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 5 {
		return 0
	}
	readBuffer.SkipHead()
	buffer, _ := readBuffer.Read(8)
	if len(buffer) != 8 {
		return 0
	}
	return BytesToFloat64(buffer)
}

func JceSectionFloat32FromBytes(readBuffer *JceReader) float32 {
	return JceSectionType4FromBytes(readBuffer)
}

func JceSectionFloat64FromBytes(readBuffer *JceReader) float64 {
	return JceSectionType5FromBytes(readBuffer)
}

func JceSectionType6FromBytes(readBuffer *JceReader) string { //jceType=6 string1
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 6 {
		return ""
	}
	readBuffer.SkipHead()
	lengthBytes, _ := readBuffer.Read(1)
	if len(lengthBytes) != 1 {
		return ""
	}
	length := uint64(lengthBytes[0])
	data, _ := readBuffer.Read(length)
	if uint64(len(data)) != length {
		return ""
	}
	return string(data)
}

func JceSectionType7FromBytes(readBuffer *JceReader) string { //jceType=7 string4
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 7 {
		return ""
	}
	readBuffer.SkipHead()
	lengthBytes, _ := readBuffer.Read(4)
	if len(lengthBytes) != 4 {
		return ""
	}
	length := uint64(BytesToInt32(lengthBytes))
	data, _ := readBuffer.Read(length)
	if uint64(len(data)) != length {
		return ""
	}
	return string(data)
}

func JceSectionStringFromBytes(readBuffer *JceReader) string {
	_, jceType, _ := readBuffer.ReadHead()
	if jceType == 6 {
		return JceSectionType6FromBytes(readBuffer)
	} else if jceType == 7 {
		return JceSectionType7FromBytes(readBuffer)
	}
	return ""
}

func JceSectionType8FromBytes(readBuffer *JceReader) (string, map[interface{}]interface{}) { //jceType=8 map
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 8 {
		return "", nil
	}
	readBuffer.SkipHead()
	length := JceSectionInt32FromBytes(readBuffer)
	if length == 0 {
		return MAP, nil
	}
	key := JceSectionStringFromBytes(readBuffer)
	if key == "" {
		return "", nil
	}
	returnMap := make(map[interface{}]interface{})
	_, jceType, _ = readBuffer.ReadHead()
	var value interface{}
	var mapType string
	if JceTypeIsString(jceType) {
		mapType = MAPStrStr
		value = JceSectionStringFromBytes(readBuffer)
	} else if JceTypeIsBytes(jceType) {
		mapType = MAPStrBytes
		value = JceSectionBytesFromBytes(readBuffer)
	} else if JceTypeIsMap(jceType) {
		value = JceSectionMapStrBytesFromBytes(readBuffer)
		if value == nil {
			return "", nil
		}
		mapType = MAPStrMAPStrBytes
	} else {
		return "", nil
	}
	returnMap[key] = value
	for i := int32(1); i < length; i++ {
		key = JceSectionStringFromBytes(readBuffer)
		if key == "" {
			return "", nil
		}
		_, jceType, _ = readBuffer.ReadHead()
		if JceTypeIsString(jceType) {
			if mapType != MAPStrStr {
				return "", nil
			}
			value = JceSectionStringFromBytes(readBuffer)
		} else if JceTypeIsBytes(jceType) {
			if mapType != MAPStrBytes {
				return "", nil
			}
			value = JceSectionBytesFromBytes(readBuffer)
		} else if JceTypeIsMap(jceType) {
			if mapType != MAPStrMAPStrBytes {
				return "", nil
			}
			value = JceSectionMapStrBytesFromBytes(readBuffer)
			if value == nil {
				return "", nil
			}
		} else {
			return "", nil
		}
		returnMap[key] = value
	}
	return mapType, returnMap
}

func JceSectionMapStrStrFromBytes(readBuffer *JceReader) map[string]string {
	mapType, interfaceMap := JceSectionType8FromBytes(readBuffer)
	if mapType != MAPStrStr {
		if mapType == MAP {
			return make(map[string]string)
		}
		return nil
	}
	returnMap := make(map[string]string)
	for k, v := range interfaceMap {
		returnMap[k.(string)] = v.(string)
	}
	return returnMap
}

func JceSectionMapStrBytesFromBytes(readBuffer *JceReader) map[string][]byte {
	mapType, interfaceMap := JceSectionType8FromBytes(readBuffer)
	if mapType != MAPStrBytes {
		if mapType == MAP {
			return make(map[string][]byte)
		}
		return nil
	}
	returnMap := make(map[string][]byte)
	for k, v := range interfaceMap {
		returnMap[k.(string)] = v.([]byte)
	}
	return returnMap
}

func JceSectionMapStrMapStrBytesFromBytes(readBuffer *JceReader) map[string]map[string][]byte {
	mapType, interfaceMap := JceSectionType8FromBytes(readBuffer)
	if mapType != MAPStrMAPStrBytes {
		if mapType == MAP {
			return make(map[string]map[string][]byte)
		}
		return nil
	}
	returnMap := make(map[string]map[string][]byte)
	for k, v := range interfaceMap {
		returnMap[k.(string)] = v.(map[string][]byte)
	}
	return returnMap
}

func JceSectionType9FromBytes(readBuffer *JceReader) (string, []interface{}) { //jceType=9 list
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 9 {
		return "", nil
	}
	readBuffer.SkipHead()
	length := JceSectionInt32FromBytes(readBuffer)
	if length == 0 {
		return LIST, nil
	}
	var returnList []interface{}
	_, jceType, _ = readBuffer.ReadHead()
	var value interface{}
	var listType string
	if JceTypeIsInt64(jceType) {
		value = JceSectionInt64FromBytes(readBuffer)
		listType = LISTInt64
	} else if JceTypeIsBytes(jceType) {
		value = JceSectionBytesFromBytes(readBuffer)
		listType = LISTBytes
	} else if JceTypeIsInt8(jceType) {
		value = JceSectionByteFromBytes(readBuffer)
		listType = BYTES
	} else {
		return "", nil
	}
	returnList = append(returnList, value)
	for i := int32(1); i < length; i++ {
		_, jceType, _ = readBuffer.ReadHead()
		if JceTypeIsInt64(jceType) {
			if listType != LISTInt64 {
				return "", nil
			}
			value = JceSectionInt64FromBytes(readBuffer)
		} else if JceTypeIsBytes(jceType) {
			if listType != LISTBytes {
				return "", nil
			}
			value = JceSectionBytesFromBytes(readBuffer)
		} else if JceTypeIsInt8(jceType) {
			if listType != BYTES {
				return "", nil
			}
			value = JceSectionByteFromBytes(readBuffer)
		} else {
			return "", nil
		}
		returnList = append(returnList, value)
	}
	return listType, returnList
}

func JceSectionListInt64FromBytes(readBuffer *JceReader) []int64 {
	listType, interfaceList := JceSectionType9FromBytes(readBuffer)
	if listType != LISTInt64 {
		if listType == LIST {
			return []int64{}
		}
		return nil
	}
	var returnList []int64
	for _, v := range interfaceList {
		returnList = append(returnList, v.(int64))
	}
	return returnList
}

func JceSectionListBytesFromBytes(readBuffer *JceReader) [][]byte {
	listType, interfaceList := JceSectionType9FromBytes(readBuffer)
	if listType != LISTBytes {
		if listType == LIST {
			return [][]byte{}
		}
		return nil
	}
	var returnList [][]byte
	for _, v := range interfaceList {
		returnList = append(returnList, v.([]byte))
	}
	return returnList
}

func JceSectionType13FromBytes(readBuffer *JceReader) []byte { //jceType=13 []byte
	_, jceType, _ := readBuffer.ReadHead()
	if jceType != 13 {
		return nil
	}
	readBuffer.SkipHead()
	readBuffer.SkipHead()
	length := uint64(JceSectionInt32FromBytes(readBuffer))
	fmt.Println(length)
	if length == 0 {
		return nil
	}
	data, _ := readBuffer.Read(length)
	if uint64(len(data)) != length {
		return nil
	}
	return data
}

func JceSectionBytesFromBytes(readBuffer *JceReader) []byte {
	_, jceType, _ := readBuffer.ReadHead()
	fmt.Println(jceType)
	if jceType == 9 {
		listType, interfaceList := JceSectionType9FromBytes(readBuffer)
		if listType != BYTES {
			if listType == LIST {
				return []byte{}
			}
			return nil
		}
		var returnList []byte
		for _, v := range interfaceList {
			returnList = append(returnList, v.(byte))
		}
		return returnList
	}
	return JceSectionType13FromBytes(readBuffer)
}

func (jceSection *JceSection) Decode(readBuffer *JceReader) (uint8, error) {
	jceId, jceType, err := readBuffer.ReadHead()
	if err != nil {
		return 0, err
	}
	//fmt.Println(jceId, jceType)
	var jceSectionType string
	var jceSectionData interface{}
	if jceType == 0 {
		jceSectionData = JceSectionType0FromBytes(readBuffer)
		jceSectionType = INT8
	} else if jceType == 1 {
		jceSectionData = JceSectionType1FromBytes(readBuffer)
		jceSectionType = INT16
	} else if jceType == 2 {
		jceSectionData = JceSectionType2FromBytes(readBuffer)
		jceSectionType = INT32
	} else if jceType == 3 {
		jceSectionData = JceSectionType3FromBytes(readBuffer)
		jceSectionType = INT64
	} else if jceType == 4 {
		jceSectionData = JceSectionType4FromBytes(readBuffer)
		jceSectionType = FLOAT32
	} else if jceType == 5 {
		jceSectionData = JceSectionType5FromBytes(readBuffer)
		jceSectionType = FLOAT64
	} else if jceType == 6 || jceType == 7 {
		jceSectionData = JceSectionStringFromBytes(readBuffer)
		jceSectionType = STRING
	} else if jceType == 8 {
		jceSectionType, jceSectionData = JceSectionType8FromBytes(readBuffer)
	} else if jceType == 9 {
		jceSectionType, jceSectionData = JceSectionType9FromBytes(readBuffer)
		if jceSectionType == BYTES {
			var data []byte
			for _, byteData := range jceSectionData.([]interface{}) {
				data = append(data, byteData.(byte))
			}
			jceSectionData = data
		}
	} else if jceType == 10 {
		jceSectionData = NewJceStruct()
		jceSectionData.(*JceStruct).structMap = jceStructFromBytes(readBuffer)
		jceSectionType = STRUCT
	} else if jceType == 11 {
		readBuffer.SkipHead()
		jceSectionType = STRUCT_END
	} else if jceType == 12 {
		readBuffer.SkipHead()
		jceSectionType = ZERO_TAG
	} else if jceType == 13 {
		jceSectionData = JceSectionType13FromBytes(readBuffer)
		jceSectionType = BYTES
	} else {
		return 0, ErrorUnknowJceTypeErr
	}
	*jceSection = JceSection{JceType: jceSectionType, Data: jceSectionData}
	return jceId, nil
}
