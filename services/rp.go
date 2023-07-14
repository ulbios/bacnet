package services

import (
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/objects"
	"github.com/ulbios/bacnet/plumbing"
)

// UnconfirmedIAm is a BACnet message.
type ConfirmedReadProperty struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type ConfirmedReadPropertyDec struct {
	ObjectType uint16
	InstanceId uint32
	PropertyId uint8
}

func ConfirmedReadPropertyObjects(objectType uint16, instN uint32, propertyId uint8) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 2)

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.EncPropertyIdentifier(true, 1, propertyId)

	return objs
}

func NewConfirmedReadProperty(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ConfirmedReadProperty {
	c := &ConfirmedReadProperty{
		BVLC: bvlc,
		NPDU: npdu,
		// TODO: Consider to implement parameter struct to an argment of New functions.
		APDU: plumbing.NewAPDU(plumbing.ConfirmedReq, ServiceConfirmedReadProperty, ConfirmedReadPropertyObjects(
			objects.ObjectTypeAnalogOutput, 1, objects.PropertyIdPresentValue)),
	}
	c.SetLength()

	return c
}

func (c *ConfirmedReadProperty) UnmarshalBinary(b []byte) error {
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

func (c *ConfirmedReadProperty) MarshalBinary() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c *ConfirmedReadProperty) MarshalTo(b []byte) error {
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

func (c *ConfirmedReadProperty) MarshalLen() int {
	l := c.BVLC.MarshalLen()
	l += c.NPDU.MarshalLen()
	l += c.APDU.MarshalLen()

	return l
}

func (c *ConfirmedReadProperty) SetLength() {
	c.BVLC.Length = uint16(c.MarshalLen())
}

func (c *ConfirmedReadProperty) Decode() (ConfirmedReadPropertyDec, error) {
	decCRP := ConfirmedReadPropertyDec{}

	if len(c.APDU.Objects) != 2 {
		return decCRP, common.ErrWrongObjectCount
	}

	for i, obj := range c.APDU.Objects {
		switch i {
		case 0:
			objId, err := objects.DecObjectIdentifier(obj)
			if err != nil {
				return decCRP, err
			}
			decCRP.ObjectType = objId.ObjectType
			decCRP.InstanceId = objId.InstanceNumber
		case 1:
			propId, err := objects.DecPropertyIdentifier(obj)
			if err != nil {
				return decCRP, err
			}
			decCRP.PropertyId = propId
		}
	}

	return decCRP, nil
}
