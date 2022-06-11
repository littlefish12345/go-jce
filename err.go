package gojce

import "errors"

var (
	ErrorUnknowJceTypeErr      = errors.New("error: Unknow JceType")
	ErrorJceStructDoesNotMatch = errors.New("error: JceStruct does not match")
)
