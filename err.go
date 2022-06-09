package gojce

import "errors"

var (
	ErrorUnknowJceTypeErr    = errors.New("error: Unknow JceType")
	ErrorJceMapInnerTypeErr  = errors.New("error: Error JceMap inner type")
	ErrorJceListInnerTypeErr = errors.New("error: Error JceList inner type")
)
