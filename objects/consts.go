package objects

// Tag number
const (
	TagZero uint8 = iota
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

const (
	ObjectTypeAnalogOutput uint16 = 1
	ObjectTypeDevice       uint16 = 8
)

const (
	PropertyIdPresentValue uint8 = 85
)

const (
	ErrorClassObject uint8 = 1

	ErrorCodeUnknownObject uint8 = 31
)
