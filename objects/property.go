package objects

import (
	"github.com/ulbios/bacnet/common"
)

func DecPropertyIdentifier(rawPayload APDUPayload) (uint8, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return 0, common.ErrWrongPayload
	}

	switch rawObject.TagClass {
	case true:
		if rawObject.Length != 1 {
			return 0, common.ErrWrongStructure
		}
	case false:
		if rawObject.Length != 1 || !rawObject.TagClass {
			return 0, common.ErrWrongStructure
		}
	}

	return rawObject.Data[0], nil
}

func EncPropertyIdentifier(contextTag bool, tagN uint8, propId uint8) *Object {
	newObj := Object{}
	data := make([]byte, 1)
	data[0] = propId

	newObj.TagNumber = tagN
	newObj.TagClass = contextTag
	newObj.Data = data
	newObj.Length = uint8(len(data))

	return &newObj
}
