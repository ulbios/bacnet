package services

import (
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/objects"
	"github.com/ulbios/bacnet/plumbing"
)

// UnconfirmedIAm is a BACnet message.
type ComplexACK struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type ComplexACKDec struct {
	ObjectType   uint16
	InstanceId   uint32
	PropertyId   uint8
	PresentValue float32
}

func ComplexACKObjects(objectType uint16, instN uint32, propertyId uint8, value float32) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 5)

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.EncPropertyIdentifier(true, 1, propertyId)
	objs[2] = objects.EncOpeningTag(3)
	objs[3] = objects.EncReal(value)
	objs[4] = objects.EncClosingTag(3)

	return objs
}

func NewComplexACK(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ComplexACK {
	c := &ComplexACK{
		BVLC: bvlc,
		NPDU: npdu,
		// TODO: Consider to implement parameter struct to an argment of New functions.
		APDU: plumbing.NewAPDU(plumbing.ComplexAck, ServiceConfirmedReadProperty, ComplexACKObjects(
			objects.ObjectTypeAnalogOutput, 1, objects.PropertyIdPresentValue, 0)),
	}
	c.SetLength()

	return c
}

func (c *ComplexACK) UnmarshalBinary(b []byte) error {
	if l := len(b); l < c.MarshalLen() {
		return common.ErrTooShortToParse
	}

	var offset int = 0
	if err := c.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}

	return nil
}

func (c *ComplexACK) MarshalBinary() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c *ComplexACK) MarshalTo(b []byte) error {
	if len(b) < c.MarshalLen() {
		return common.ErrTooShortToMarshalBinary
	}
	var offset = 0
	if err := c.BVLC.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.MarshalTo(b[offset:]); err != nil {
		return err
	}

	return nil
}

func (c *ComplexACK) MarshalLen() int {
	l := c.BVLC.MarshalLen()
	l += c.NPDU.MarshalLen()
	l += c.APDU.MarshalLen()

	return l
}

func (u *ComplexACK) SetLength() {
	u.BVLC.Length = uint16(u.MarshalLen())
}

func (c *ComplexACK) Decode() (ComplexACKDec, error) {
	decCACK := ComplexACKDec{}

	if len(c.APDU.Objects) != 3 {
		return decCACK, common.ErrWrongObjectCount
	}

	for i, obj := range c.APDU.Objects {
		switch i {
		case 0:
			objId, err := objects.DecObjectIdentifier(obj)
			if err != nil {
				return decCACK, err
			}
			decCACK.ObjectType = objId.ObjectType
			decCACK.InstanceId = objId.InstanceNumber
		case 1:
			propId, err := objects.DecPropertyIdentifier(obj)
			if err != nil {
				return decCACK, err
			}
			decCACK.PropertyId = propId
		case 3:
			value, err := objects.DecReal(obj)
			if err != nil {
				return decCACK, err
			}
			decCACK.PresentValue = value
		}
	}

	return decCACK, nil
}
