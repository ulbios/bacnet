package services

import (
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/objects"
	"github.com/ulbios/bacnet/plumbing"
)

// UnconfirmedIAm is a BACnet message.
type Error struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type ErrorDec struct {
	ErrorClass uint8
	ErrorCode  uint8
}

// IAmObjects creates an instance of UnconfirmedIAm objects.
func ErrorObjects(errClass, errCode uint8) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 2)

	objs[0] = objects.EncEnumerated(errClass)
	objs[1] = objects.EncEnumerated(errCode)

	return objs
}

// NewUnconfirmedIAm creates a UnconfirmedIam.
func NewError(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *Error {
	e := &Error{
		BVLC: bvlc,
		NPDU: npdu,
		// TODO: Consider to implement parameter struct to an argment of New functions.
		APDU: plumbing.NewAPDU(plumbing.Error, ServiceConfirmedReadProperty, ErrorObjects(1, 31)),
	}
	e.SetLength()

	return e
}

// UnmarshalBinary sets the values retrieved from byte sequence in a UnconfirmedIAm frame.
func (e *Error) UnmarshalBinary(b []byte) error {
	if l := len(b); l < e.MarshalLen() {
		return common.ErrTooShortToParse
	}

	var offset int = 0
	if err := e.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += e.BVLC.MarshalLen()

	if err := e.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += e.NPDU.MarshalLen()

	if err := e.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}

	return nil
}

// MarshalBinary returns the byte sequence generated from a UnconfirmedIAm instance.
func (e *Error) MarshalBinary() ([]byte, error) {
	b := make([]byte, e.MarshalLen())
	if err := e.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (e *Error) MarshalTo(b []byte) error {
	if len(b) < e.MarshalLen() {
		return common.ErrTooShortToMarshalBinary
	}
	var offset = 0
	if err := e.BVLC.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += e.BVLC.MarshalLen()

	if err := e.NPDU.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += e.NPDU.MarshalLen()

	if err := e.APDU.MarshalTo(b[offset:]); err != nil {
		return err
	}

	return nil
}

// MarshalLen returns the serial length of UnconfirmedIAm.
func (e *Error) MarshalLen() int {
	l := e.BVLC.MarshalLen()
	l += e.NPDU.MarshalLen()
	l += e.APDU.MarshalLen()

	return l
}

// SetLength sets the length in Length field.
func (e *Error) SetLength() {
	e.BVLC.Length = uint16(e.MarshalLen())
}

func (e *Error) Decode() (ErrorDec, error) {
	decErr := ErrorDec{}

	if len(e.APDU.Objects) != 2 {
		return decErr, common.ErrWrongObjectCount
	}

	for i, obj := range e.APDU.Objects {
		switch i {
		case 0:
			errClass, err := objects.DecEnumerated(obj)
			if err != nil {
				return decErr, err
			}
			decErr.ErrorClass = uint8(errClass)
		case 1:
			errCode, err := objects.DecEnumerated(obj)
			if err != nil {
				return decErr, err
			}
			decErr.ErrorCode = uint8(errCode)
		}
	}

	return decErr, nil
}
