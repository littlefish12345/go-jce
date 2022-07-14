package gojce

import (
	"fmt"
	"reflect"
	"strconv"
)

type JceStruct struct {
	structMap map[uint8]*JceSection
}

func (jceStruct *JceStruct) SetSection(jceId uint8, jceSection JceSection) {
	jceStruct.structMap[jceId] = &jceSection
}

func (jceStruct *JceStruct) CopyTemplate() *JceStruct {
	newJceStruct := NewJceStruct()
	for k, v := range jceStruct.structMap {
		newJceStruct.structMap[k] = v
	}
	return newJceStruct
}

func NewJceStruct() *JceStruct {
	return &JceStruct{make(map[uint8]*JceSection)}
}

func Marshal(vStruct interface{}) (*JceStruct, error) {
	jceStruct := NewJceStruct()
	typeOfvStruct := reflect.TypeOf(vStruct)
	valueOfvStruct := reflect.ValueOf(vStruct)
	var jceSectionType string
	var jceSectionData interface{}
	var jceId uint64
	var err error
	var fieldType reflect.StructField
	var nameField reflect.StructField
	for i := 0; i < typeOfvStruct.NumField(); i++ {
		fieldType = typeOfvStruct.Field(i)
		nameField = typeOfvStruct.FieldByIndex(fieldType.Index)
		jceId, err = strconv.ParseUint(nameField.Tag.Get("jceId"), 10, 8)
		if err != nil {
			return nil, err
		}
		fieldTypeString := fieldType.Type.String()
		if fieldTypeString == reflect.Uint8.String() {
			jceSectionType = BYTE
			jceSectionData = byte(valueOfvStruct.FieldByName(fieldType.Name).Uint())
		} else if fieldTypeString == reflect.Int8.String() {
			jceSectionType = INT8
			jceSectionData = int8(valueOfvStruct.FieldByName(fieldType.Name).Int())
		} else if fieldTypeString == reflect.Int16.String() {
			jceSectionType = INT16
			jceSectionData = int16(valueOfvStruct.FieldByName(fieldType.Name).Int())
		} else if fieldTypeString == reflect.Int32.String() {
			jceSectionType = INT32
			jceSectionData = int32(valueOfvStruct.FieldByName(fieldType.Name).Int())
		} else if fieldTypeString == reflect.Int64.String() {
			jceSectionType = INT64
			jceSectionData = int64(valueOfvStruct.FieldByName(fieldType.Name).Int())
		} else if fieldTypeString == reflect.Float32.String() {
			jceSectionType = FLOAT32
			jceSectionData = float32(valueOfvStruct.FieldByName(fieldType.Name).Float())
		} else if fieldTypeString == reflect.Float64.String() {
			jceSectionType = FLOAT64
			jceSectionData = float64(valueOfvStruct.FieldByName(fieldType.Name).Float())
		} else if fieldTypeString == reflect.String.String() {
			jceSectionType = STRING
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).String()
		} else if fieldTypeString == "[]uint8" {
			jceSectionType = BYTES
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Interface()
		} else if fieldTypeString == reflect.Bool.String() {
			jceSectionType = BOOL
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Bool()
		} else if fieldTypeString == "map[string]string" {
			jceSectionType = MAPStrStr
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Interface()
		} else if fieldTypeString == "map[string][]uint8" {
			jceSectionType = MAPStrBytes
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Interface()
		} else if fieldTypeString == "map[string]map[string]uint8" {
			jceSectionType = MAPStrMAPStrBytes
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Interface()
		} else if fieldTypeString == "[]int64" {
			jceSectionType = LISTInt64
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Interface()
		} else if fieldTypeString == "[][]uint8" {
			jceSectionType = LISTBytes
			jceSectionData = valueOfvStruct.FieldByName(fieldType.Name).Interface()
		} else {
			fmt.Println(fieldTypeString)
			return nil, ErrorUnknowJceTypeErr
		}
		jceStruct.SetSection(uint8(jceId), JceSection{JceType: jceSectionType, Data: jceSectionData})
	}
	return jceStruct, nil
}

func Unmarshal(data []byte, vStruct interface{}) error {
	jceStruct := NewJceStruct()
	err := jceStruct.Decode(data)
	if err != nil {
		return err
	}
	err = Bind(jceStruct, vStruct)
	return err
}

func Bind(jceStruct *JceStruct, vStruct interface{}) error {
	typeOfvStruct := reflect.TypeOf(vStruct).Elem()
	valueOfvStruct := reflect.ValueOf(vStruct)
	var err error
	var jceId uint64
	var jceSection *JceSection
	var ok bool
	var fieldType reflect.StructField
	var nameField reflect.StructField
	for i := 0; i < typeOfvStruct.NumField(); i++ {
		fieldType = typeOfvStruct.Field(i)
		nameField = typeOfvStruct.FieldByIndex(fieldType.Index)
		jceId, err = strconv.ParseUint(nameField.Tag.Get("jceId"), 10, 8)
		if err != nil {
			return err
		}
		fieldTypeString := fieldType.Type.String()
		if jceSection, ok = jceStruct.structMap[uint8(jceId)]; !ok {
			continue
		}
		//fmt.Println(jceSection.JceType, jceId, fieldTypeString)
		if fieldTypeString == reflect.Uint8.String() {
			if jceSection.JceType == INT8 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetUint(uint64(jceSection.Data.(int8)))
			} else if jceSection.JceType == ZERO_TAG {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetUint(0)
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == reflect.Int8.String() {
			if jceSection.JceType == INT8 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(int64(jceSection.Data.(int8)))
			} else if jceSection.JceType == ZERO_TAG {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(0)
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == reflect.Int16.String() {
			if jceSection.JceType == INT8 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(int64(jceSection.Data.(int8)))
			} else if jceSection.JceType == INT16 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(int64(jceSection.Data.(int16)))
			} else if jceSection.JceType == ZERO_TAG {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(0)
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == reflect.Int32.String() {
			if jceSection.JceType == INT8 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(int64(jceSection.Data.(int8)))
			} else if jceSection.JceType == INT16 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(int64(jceSection.Data.(int16)))
			} else if jceSection.JceType == INT32 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(int64(jceSection.Data.(int32)))
			} else if jceSection.JceType == ZERO_TAG {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(0)
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == reflect.Int64.String() {
			if jceSection.JceType == INT8 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(int64(jceSection.Data.(int8)))
			} else if jceSection.JceType == INT16 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(int64(jceSection.Data.(int16)))
			} else if jceSection.JceType == INT32 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(int64(jceSection.Data.(int32)))
			} else if jceSection.JceType == INT64 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(jceSection.Data.(int64))
			} else if jceSection.JceType == ZERO_TAG {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetInt(0)
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == reflect.Float32.String() {
			if jceSection.JceType == FLOAT32 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetFloat(float64(jceSection.Data.(float32)))
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == reflect.Float64.String() {
			if jceSection.JceType == FLOAT64 {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetFloat(jceSection.Data.(float64))
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == reflect.String.String() {
			if jceSection.JceType == STRING {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetString(jceSection.Data.(string))
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == "[]uint8" {
			if jceSection.JceType == BYTES {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.ValueOf(jceSection.Data.([]uint8)))
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == reflect.Bool.String() {
			if jceSection.JceType == INT8 {
				if jceSection.Data.(int8) == 0x01 {
					valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetBool(true)
				} else {
					return ErrorJceStructDoesNotMatch
				}
			} else if jceSection.JceType == ZERO_TAG {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).SetBool(false)
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == "map[string]string" {
			if jceSection.JceType == MAPStrStr {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.ValueOf(jceSection.Data.(map[string]string)))
			} else if jceSection.JceType == MAP {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.Zero(fieldType.Type))
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == "map[string][]uint8" {
			if jceSection.JceType == MAPStrBytes {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.ValueOf(jceSection.Data.(map[string][]uint8)))
			} else if jceSection.JceType == MAP {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.Zero(fieldType.Type))
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == "map[string]map[string][]uint8" {
			if jceSection.JceType == MAPStrMAPStrBytes {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.ValueOf(jceSection.Data.(map[string]map[string][]uint8)))
			} else if jceSection.JceType == MAP {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.Zero(fieldType.Type))
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == "[]int64" {
			if jceSection.JceType == LISTBytes {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.ValueOf(jceSection.Data.([]int64)))
			} else if jceSection.JceType == LIST {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Elem().Set(reflect.Zero(fieldType.Type))
			} else {
				return ErrorJceStructDoesNotMatch
			}
		} else if fieldTypeString == "[][]uint8" {
			if jceSection.JceType == LISTBytes {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.ValueOf(jceSection.Data.([][]uint8)))
			} else if jceSection.JceType == LIST {
				valueOfvStruct.Elem().FieldByIndex(fieldType.Index).Set(reflect.Zero(fieldType.Type))
			} else {
				return ErrorJceStructDoesNotMatch
			}
		}
	}
	return nil
}
