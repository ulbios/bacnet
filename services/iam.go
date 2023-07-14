package services

import (
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/objects"
	"github.com/ulbios/bacnet/plumbing"
)

// UnconfirmedIAm is a BACnet message.
type UnconfirmedIAm struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type UnconfirmedIAmDec struct {
	DeviceId              uint32
	MaxAPDULength         uint16
	SegmentationSupported uint8
	VendorId              uint16
}

// IAmObjects creates an instance of UnconfirmedIAm objects.
func IAmObjects(insNum uint32, acceptedSize uint16, supportedSeg uint8, vendorID uint16) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 4)

	objs[0] = objects.EncObjectIdentifier(false, objects.TagBACnetObjectIdentifier, objects.ObjectTypeDevice, 321)
	objs[1] = objects.EncUnsignedInteger16(acceptedSize)
	objs[2] = objects.EncEnumerated(supportedSeg)
	objs[3] = objects.EncUnsignedInteger16(vendorID)

	return objs
}

// NewUnconfirmedIAm creates a UnconfirmedIam.
func NewUnconfirmedIAm(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *UnconfirmedIAm {
	u := &UnconfirmedIAm{
		BVLC: bvlc,
		NPDU: npdu,
		// TODO: Consider to implement parameter struct to an argment of New functions.
		APDU: plumbing.NewAPDU(plumbing.UnConfirmedReq, ServiceUnconfirmedIAm, IAmObjects(1, 1024, 0, 1)),
	}
	u.SetLength()

	return u
}

// UnmarshalBinary sets the values retrieved from byte sequence in a UnconfirmedIAm frame.
func (u *UnconfirmedIAm) UnmarshalBinary(b []byte) error {
	if l := len(b); l < u.MarshalLen() {
		return common.ErrTooShortToParse
	}

	var offset int = 0
	if err := u.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += u.BVLC.MarshalLen()

	if err := u.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += u.NPDU.MarshalLen()

	if err := u.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}

	return nil
}

// MarshalBinary returns the byte sequence generated from a UnconfirmedIAm instance.
func (u *UnconfirmedIAm) MarshalBinary() ([]byte, error) {
	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (u *UnconfirmedIAm) MarshalTo(b []byte) error {
	if len(b) < u.MarshalLen() {
		return common.ErrTooShortToMarshalBinary
	}
	var offset = 0
	if err := u.BVLC.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += u.BVLC.MarshalLen()

	if err := u.NPDU.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += u.NPDU.MarshalLen()

	if err := u.APDU.MarshalTo(b[offset:]); err != nil {
		return err
	}

	return nil
}

// MarshalLen returns the serial length of UnconfirmedIAm.
func (u *UnconfirmedIAm) MarshalLen() int {
	l := u.BVLC.MarshalLen()
	l += u.NPDU.MarshalLen()
	l += u.APDU.MarshalLen()

	return l
}

// SetLength sets the length in Length field.
func (u *UnconfirmedIAm) SetLength() {
	u.BVLC.Length = uint16(u.MarshalLen())
}

func (u *UnconfirmedIAm) Decode() (UnconfirmedIAmDec, error) {
	decIAm := UnconfirmedIAmDec{}

	if len(u.APDU.Objects) != 4 {
		return decIAm, common.ErrWrongObjectCount
	}

	for i, obj := range u.APDU.Objects {
		switch i {
		case 0:
			objId, err := objects.DecObjectIdentifier(obj)
			if err != nil {
				return decIAm, err
			}
			decIAm.DeviceId = objId.InstanceNumber
		case 1:
			maxLen, err := objects.DecUnisgnedInteger(obj)
			if err != nil {
				return decIAm, err
			}
			decIAm.MaxAPDULength = uint16(maxLen)
		case 2:
			segSupport, err := objects.DecEnumerated(obj)
			if err != nil {
				return decIAm, err
			}
			decIAm.SegmentationSupported = uint8(segSupport)
		case 3:
			vendorId, err := objects.DecUnisgnedInteger(obj)
			if err != nil {
				return decIAm, err
			}
			decIAm.VendorId = uint16(vendorId)
		}
	}

	return decIAm, nil
}
