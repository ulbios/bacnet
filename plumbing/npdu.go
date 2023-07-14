package plumbing

import (
	"encoding/binary"

	"github.com/ulbios/bacnet/common"
)

// NPDU is a Network Protocol Data Units.
type NPDU struct {
	Version uint8
	Control uint8
	DNET    uint16
	DLEN    uint8
	Hop     uint8
}

// NewNPDU creates a NPDU.
func NewNPDU(nsduContain bool, dstSpecifier bool, srcSpecifier bool, expectingReply bool) *NPDU {
	n := &NPDU{
		Version: 1,
	}
	n.SetControlFlags(nsduContain, dstSpecifier, srcSpecifier, expectingReply)
	return n
}

// SetControlFlags sets control flags to NPDU.
func (n *NPDU) SetControlFlags(nsduContain bool, dstSpecifier bool, srcSpecifier bool, expectingReply bool) {
	n.Control = uint8(
		common.BoolToInt(nsduContain)<<7 | common.BoolToInt(dstSpecifier)<<5 |
			common.BoolToInt(srcSpecifier)<<3 | common.BoolToInt(expectingReply)<<2,
	)
}

// UnmarshalBinary sets the values retrieved from byte sequence in a NPDU frame.
func (n *NPDU) UnmarshalBinary(b []byte) error {
	if l := len(b); l < n.MarshalLen() {
		return common.ErrTooShortToParse
	}
	n.Version = b[0]
	n.Control = b[1]
	if flagDNET := n.Control & 0x20 >> 5; flagDNET == 1 {
		n.DNET = binary.BigEndian.Uint16(b[2:4])
		n.DLEN = b[4]
		n.Hop = b[5]
	}

	return nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (n *NPDU) MarshalTo(b []byte) error {
	if len(b) < n.MarshalLen() {
		return common.ErrTooShortToMarshalBinary
	}
	b[0] = n.Version
	b[1] = n.Control
	if flagDNET := n.Control & 0x20 >> 5; flagDNET == 1 {
		binary.BigEndian.PutUint16(b[2:4], n.DNET)
		b[4] = n.DLEN
		b[5] = n.Hop
	}
	return nil
}

const npduLenMin = 2

// MarshalLen returns the serial length of NPDU.
func (n *NPDU) MarshalLen() int {
	flagDNET := n.Control & 0x20 >> 5
	if flagDNET == 1 {
		return npduLenMin + 4
	}
	return npduLenMin
}
