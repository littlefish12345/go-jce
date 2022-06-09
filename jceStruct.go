package gojce

import (
	"reflect"
	"strconv"
)

type JceStruct struct {
	structMap map[uint8]JceSection
}

func (jceStruct *JceStruct) SetSection(jceId uint8, jceSection JceSection) {
	jceStruct.structMap[jceId] = jceSection
}

func (jceStruct *JceStruct) CopyTemplate() *JceStruct {
	newJceStruct := newJceStruct()
	for k, v := range jceStruct.structMap {
		newJceStruct.structMap[k] = v
	}
	return newJceStruct
}

func newJceStruct() *JceStruct {
	return &JceStruct{make(map[uint8]JceSection)}
}

func Marshal(vStruct interface{}) (*JceStruct, error) {
	jceStruct := newJceStruct()
	typeOfvStruct := reflect.TypeOf(vStruct)
	valueOfvStruct := reflect.ValueOf(vStruct)
	var jceSectionType string
	var jceSectionData interface{}
	var jceId uint64
	var err error
	var nameField reflect.StructField
	for i := 0; i < typeOfvStruct.NumField(); i++ {
		fieldType := typeOfvStruct.Field(i)
		nameField = typeOfvStruct.FieldByIndex(fieldType.Index)
		jceId, err = strconv.ParseUint(nameField.Tag.Get("jceId"), 10, 8)
		if err != nil {
			return nil, err
		}
		if fieldType.Type.String() == reflect.Uint8.String() {
			jceSectionType = BYTE
			jceSectionData = byte(valueOfvStruct.FieldByName(fieldType.Name).Uint())
		} else if fieldType.Type.String() == reflect.Int8.String() {
			jceSectionType = INT8
			jceSectionData = int8(valueOfvStruct.FieldByName(fieldType.Name).Int())
		} else if fieldType.Type.String() == reflect.Int16.String() {
			jceSectionType = INT16
			jceSectionData = int16(valueOfvStruct.FieldByName(fieldType.Name).Int())
		} else if fieldType.Type.String() == reflect.Int32.String() {
			jceSectionType = INT32
			jceSectionData = int32(valueOfvStruct.FieldByName(fieldType.Name).Int())
		} else if fieldType.Type.String() == reflect.Int64.String() {
			jceSectionType = INT64
			jceSectionData = int64(valueOfvStruct.FieldByName(fieldType.Name).Int())
		} else if fieldType.Type.String() == reflect.Float32.String() {
			jceSectionType = FLOAT32
			jceSectionData = float32(valueOfvStruct.FieldByName(fieldType.Name).Float())
		} else if fieldType.Type.String() == reflect.Float64.String() {
			jceSectionType = FLOAT64
			jceSectionData = float64(valueOfvStruct.FieldByName(fieldType.Name).Float())
		} else if fieldType.Type.String() == reflect.String.String() {
			jceSectionType = STRING
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).String()
		} else if fieldType.Type.String() == reflect.Bool.String() {
			jceSectionType = BOOL
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Bool()
		} else if fieldType.Type.String() == reflect.Map.String() {
			if fieldType.Type.Key().String() != reflect.String.String() {
				return nil, ErrorJceMapInnerTypeErr
			}
			if fieldType.Type.Elem().String() == reflect.String.String() {
				jceSectionType = MAPStrStr
			} else if fieldType.Type.Elem().String() == reflect.Slice.String() && fieldType.Type.Elem().Elem().String() == reflect.Uint8.String() {
				jceSectionType = MAPStrBytes
			} else if fieldType.Type.Elem().String() == reflect.Map.String() {
				if fieldType.Type.Elem().Key().String() != reflect.String.String() {
					return nil, ErrorJceMapInnerTypeErr
				}
				if fieldType.Type.Elem().Elem().String() != reflect.Slice.String() || fieldType.Type.Elem().Elem().Elem().String() != reflect.Uint8.String() {
					return nil, ErrorJceMapInnerTypeErr
				}
				jceSectionType = MAPStrMAPStrBytes
			} else {
				return nil, ErrorJceMapInnerTypeErr
			}
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Interface()
		} else if fieldType.Type.String() == reflect.Slice.String() {
			if fieldType.Type.Elem().String() == reflect.Int64.String() {
				jceSectionType = LISTInt64
			} else if fieldType.Type.Elem().String() == reflect.Uint8.String() {
				jceSectionType = LISTBytes
			} else {
				return nil, ErrorJceListInnerTypeErr
			}
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Interface()
		} else {
			return nil, ErrorUnknowJceTypeErr
		}
		jceStruct.SetSection(uint8(jceId), JceSection{JceType: jceSectionType, Data: jceSectionData})
	}
	return jceStruct, nil
}
