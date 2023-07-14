package services

import (
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/plumbing"
)

// UnconfirmedIAm is a BACnet message.
type SimpleACK struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

func NewSimpleACK(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *SimpleACK {
	s := &SimpleACK{
		BVLC: bvlc,
		NPDU: npdu,
		// TODO: Consider to implement parameter struct to an argment of New functions.
		APDU: plumbing.NewAPDU(plumbing.SimpleAck, ServiceConfirmedReadProperty, nil),
	}
	s.SetLength()

	return s
}

func (s *SimpleACK) UnmarshalBinary(b []byte) error {
	if l := len(b); l < s.MarshalLen() {
		return common.ErrTooShortToParse
	}

	var offset int = 0
	if err := s.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += s.BVLC.MarshalLen()

	if err := s.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += s.NPDU.MarshalLen()

	if err := s.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}

	return nil
}

func (s *SimpleACK) MarshalBinary() ([]byte, error) {
	b := make([]byte, s.MarshalLen())
	if err := s.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *SimpleACK) MarshalTo(b []byte) error {
	if len(b) < s.MarshalLen() {
		return common.ErrTooShortToMarshalBinary
	}
	var offset = 0
	if err := s.BVLC.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += s.BVLC.MarshalLen()

	if err := s.NPDU.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += s.NPDU.MarshalLen()

	if err := s.APDU.MarshalTo(b[offset:]); err != nil {
		return err
	}

	return nil
}

func (s *SimpleACK) MarshalLen() int {
	l := s.BVLC.MarshalLen()
	l += s.NPDU.MarshalLen()
	l += s.APDU.MarshalLen()

	return l
}

func (s *SimpleACK) SetLength() {
	s.BVLC.Length = uint16(s.MarshalLen())
}
