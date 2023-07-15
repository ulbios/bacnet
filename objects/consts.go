package objects

// Tag number
const (
	TagNull uint8 = iota
	TagBoolean
	TagUnsignedInteger
	TagSignedInteger
	TagReal
	TagDouble
	TagOctetString
	TagCharacterString
	TagBitString
	TagEnumerated
	TagDate
	TagTime
	TagBACnetObjectIdentifier
	TagOpening uint8 = 0x3E
	TagClosing uint8 = 0x3F
)

// Be sure to check ../bacnet-stack/src/bacnet/bacenum.h for more!
const (
	ObjectTypeAnalogInput  uint16 = 0
	ObjectTypeAnalogOutput uint16 = 1
	ObjectTypeDevice       uint16 = 8
)

const (
	PropertyIdPresentValue uint8 = 85
)

const (
	ErrorClassObject  uint8 = 1
	ErrorClassService uint8 = 5

	ErrorCodeUnknownObject        uint8 = 31
	ErrorCodeServiceRequestDenied uint8 = 29
)
