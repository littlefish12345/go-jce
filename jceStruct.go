package gojce

import (
	"fmt"
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
