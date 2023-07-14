package objects

import (
	"encoding/binary"
	"math"

	"github.com/ulbios/bacnet/common"
)

func DecUnisgnedInteger(rawPayload APDUPayload) (uint32, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return 0, common.ErrWrongPayload
	}

	if rawObject.TagNumber != TagUnsignedInteger || rawObject.TagClass {
		return 0, common.ErrWrongStructure
	}

	switch rawObject.Length {
	case 1:
		return uint32(rawObject.Data[0]), nil
	case 2:
		return uint32(binary.BigEndian.Uint16(rawObject.Data)), nil
	case 3:
		return uint32(uint16(uint32(rawObject.Data[0])<<16) | binary.BigEndian.Uint16(rawObject.Data[1:])), nil
	case 4:
		return binary.BigEndian.Uint32(rawObject.Data), nil
	}

	return 0, common.ErrNotImplemented
}

func EncUnsignedInteger16(value uint16) *Object {
	newObj := Object{}

	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data[:], value)

	newObj.TagNumber = TagUnsignedInteger
	newObj.TagClass = false
	newObj.Data = data
	newObj.Length = uint8(len(data))

	return &newObj
}

func DecEnumerated(rawPayload APDUPayload) (uint32, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return 0, common.ErrWrongPayload
	}

	if rawObject.TagNumber != TagEnumerated || rawObject.TagClass {
		return 0, common.ErrWrongStructure
	}

	switch rawObject.Length {
	case 1:
		return uint32(rawObject.Data[0]), nil
	case 2:
		return uint32(binary.BigEndian.Uint16(rawObject.Data)), nil
	case 3:
		return uint32(uint16(uint32(rawObject.Data[0])<<16) | binary.BigEndian.Uint16(rawObject.Data[1:])), nil
	case 4:
		return binary.BigEndian.Uint32(rawObject.Data), nil
	}

	return 0, common.ErrNotImplemented
}

func EncEnumerated(value uint8) *Object {
	newObj := Object{}

	data := make([]byte, 1)
	data[0] = value

	newObj.TagNumber = TagEnumerated
	newObj.TagClass = false
	newObj.Data = data
	newObj.Length = uint8(len(data))

	return &newObj
}

func DecReal(rawPayload APDUPayload) (float32, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return 0, common.ErrWrongPayload
	}

	if rawObject.TagNumber != TagReal {
		return 0, common.ErrWrongStructure
	}

	return math.Float32frombits(binary.BigEndian.Uint32(rawObject.Data)), nil
}

func EncReal(value float32) *Object {
	newObj := Object{}

	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data[:], math.Float32bits(value))

	newObj.TagNumber = TagReal
	newObj.TagClass = false
	newObj.Data = data
	newObj.Length = uint8(len(data))

	return &newObj
}

func DecNull(rawPayload APDUPayload) (bool, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return false, common.ErrWrongPayload
	}

	if rawObject.TagNumber != TagReal {
		return false, common.ErrWrongStructure
	}

	return rawObject.TagNumber == TagNull && !rawObject.TagClass && rawObject.Length == 0, nil
}

func EncNull() *Object {
	newObj := Object{}

	newObj.TagNumber = TagNull
	newObj.TagClass = false
	newObj.Data = nil
	newObj.Length = 0

	return &newObj
}
