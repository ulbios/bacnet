package bacnet

import (
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/plumbing"
	"github.com/ulbios/bacnet/services"
)

const bacnetLenMin = 8

func combine(t, s uint8) uint16 {
	return uint16(t)<<8 | uint16(s)
}

// Parse decodes the given bytes.
func Parse(b []byte) (plumbing.BACnet, error) {
	if len(b) < bacnetLenMin {
		return nil, common.ErrTooShortToParse
	}

	var bvlc plumbing.BVLC
	var npdu plumbing.NPDU
	var bacnet plumbing.BACnet

	offset := 0

	if err := bvlc.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	offset += bvlc.MarshalLen()

	if err := npdu.UnmarshalBinary(b[offset:]); err != nil {
		return nil, err
	}
	offset += npdu.MarshalLen()

	var c uint16
	switch b[offset] >> 4 & 0xFF {
	case plumbing.UnConfirmedReq:
		c = combine(b[offset], b[offset+1])
	case plumbing.ConfirmedReq:
		c = combine(b[offset], b[offset+3]) // We need to skip the PDU flags and the InvokeID
	case plumbing.ComplexAck, plumbing.SimpleAck, plumbing.Error:
		c = combine(b[offset], 0) // We need to skip the PDU flags and the InvokeID
	}

	switch c {
	case combine(plumbing.UnConfirmedReq<<4, services.ServiceUnconfirmedWhoIs):
		bacnet = services.NewUnconfirmedWhoIs(&bvlc, &npdu)
	case combine(plumbing.UnConfirmedReq<<4, services.ServiceUnconfirmedIAm):
		bacnet = services.NewUnconfirmedIAm(&bvlc, &npdu)
	case combine(plumbing.ConfirmedReq<<4, services.ServiceConfirmedReadProperty):
		bacnet = services.NewConfirmedReadProperty(&bvlc, &npdu)
	case combine(plumbing.ConfirmedReq<<4, services.ServiceConfirmedWriteProperty):
		bacnet = services.NewConfirmedWriteProperty(&bvlc, &npdu)
	case combine(plumbing.ComplexAck<<4, 0):
		bacnet = services.NewComplexACK(&bvlc, &npdu)
	case combine(plumbing.SimpleAck<<4, 0):
		bacnet = services.NewSimpleACK(&bvlc, &npdu)
	case combine(plumbing.Error<<4, 0):
		bacnet = services.NewError(&bvlc, &npdu)
	default:
		return nil, common.ErrNotImplemented
	}

	if err := bacnet.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return bacnet, nil
}
