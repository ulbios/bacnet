package objects

import (
	"github.com/ulbios/bacnet/common"
)

type APDUPayload interface {
	UnmarshalBinary([]byte) error
	MarshalBinary() ([]byte, error)
	MarshalTo([]byte) error
	MarshalLen() int
}

// Object is an object in APDU.
type Object struct {
	TagNumber uint8
	TagClass  bool
	Length    uint8
	Data      []byte
}

// NewObject creates an Object.
func NewObject(number uint8, class bool, data []byte) *Object {
	return &Object{
		TagNumber: number,
		TagClass:  class,
		Data:      data,
	}
}

const objLenMin int = 2

// UnmarshalBinary sets the values retrieved from byte sequence in a Object frame.
func (o *Object) UnmarshalBinary(b []byte) error {
	if l := len(b); l < objLenMin {
		return common.ErrTooShortToParse
	}
	o.TagNumber = b[0] >> 4
	o.TagClass = common.IntToBool(int(b[0]) & 0x8 >> 3)
	o.Length = b[0] & 0x7

	if l := len(b); l < int(o.Length) {
		return common.ErrTooShortToParse
	}

	o.Data = b[1:o.Length]

	return nil
}

// MarshalBinary returns the byte sequence generated from a Object instance.
func (o *Object) MarshalBinary() ([]byte, error) {
	b := make([]byte, o.MarshalLen())
	if err := o.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (o *Object) MarshalTo(b []byte) error {
	if len(b) < o.MarshalLen() {
		return common.ErrTooShortToMarshalBinary
	}
	b[0] = o.TagNumber<<4 | uint8(common.BoolToInt(o.TagClass))<<3 | o.Length
	if o.Length > 0 {
		copy(b[1:o.Length+1], o.Data)
	}
	return nil
}

// MarshalLen returns the serial length of Object.
func (o *Object) MarshalLen() int {
	return 1 + int(o.Length)
}
