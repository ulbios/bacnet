package services

import (
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/plumbing"
)

// UnconfirmedWhoIs is a BACnet message.
type UnconfirmedWhoIs struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

// NewUnconfirmedWhoIs creates a UnconfirmedWhoIs.
func NewUnconfirmedWhoIs(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *UnconfirmedWhoIs {
	u := &UnconfirmedWhoIs{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: plumbing.NewAPDU(plumbing.UnConfirmedReq, ServiceUnconfirmedWhoIs, nil),
	}
	u.SetLength()
	return u
}

// UnmarshalBinary sets the values retrieved from byte sequence in a UnconfirmedWhoIs frame.
func (u *UnconfirmedWhoIs) UnmarshalBinary(b []byte) error {
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

// MarshalBinary returns the byte sequence generated from a UnconfirmedWhoIs instance.
func (u *UnconfirmedWhoIs) MarshalBinary() ([]byte, error) {
	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (u *UnconfirmedWhoIs) MarshalTo(b []byte) error {
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

// MarshalLen returns the serial length of UnconfirmedWhoIs.
func (u *UnconfirmedWhoIs) MarshalLen() int {
	l := u.BVLC.MarshalLen()
	l += u.NPDU.MarshalLen()
	l += u.APDU.MarshalLen()

	return l
}

// SetLength sets the length in Length field.
func (u *UnconfirmedWhoIs) SetLength() {
	u.BVLC.Length = uint16(u.MarshalLen())
}
