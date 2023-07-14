package plumbing

import (
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/objects"
)

// APDU is a Application protocol DAta Units.
type APDU struct {
	Type     uint8
	Flags    uint8
	MaxSeg   uint8
	MaxSize  uint8
	InvokeID uint8
	Service  uint8
	Objects  []objects.APDUPayload
}

// NewAPDU creates an APDU.
func NewAPDU(t, s uint8, objs []objects.APDUPayload) *APDU {
	return &APDU{
		Type:    t,
		Service: s,
		Objects: objs,
	}
}

// UnmarshalBinary sets the values retrieved from byte sequence in a APDU frame.
func (a *APDU) UnmarshalBinary(b []byte) error {
	if l := len(b); l < a.MarshalLen() {
		return common.ErrTooShortToParse
	}

	a.Type = b[0] >> 4
	a.Flags = b[0] & 0x7

	var offset int = 1

	switch a.Type {
	case UnConfirmedReq:
		a.Service = b[offset]
		offset++
		if len(b) > 2 {
			objs := []objects.APDUPayload{}
			for {
				o := objects.Object{
					TagNumber: b[offset] >> 4,
					TagClass:  common.IntToBool(int(b[offset]) & 0x8 >> 3),
					Length:    b[offset] & 0x7,
				}

				o.Data = b[offset+1 : offset+int(o.Length)+1]
				objs = append(objs, &o)
				offset += int(o.Length) + 1

				if offset >= len(b) {
					break
				}
			}
			a.Objects = objs
		}
	case ConfirmedReq:
		offset++
		a.InvokeID = b[offset]
		offset++
		a.Service = b[offset]
		offset++
		if len(b) > 2 {
			objs := []objects.APDUPayload{}
			for {
				o := objects.Object{
					TagNumber: b[offset] >> 4,
					TagClass:  common.IntToBool(int(b[offset]) & 0x8 >> 3),
					Length:    b[offset] & 0x7,
				}

				// Drop tags so that they don't get in the way!
				if b[offset] == objects.TagOpening || b[offset] == objects.TagClosing {
					offset++
					if offset >= len(b) {
						break
					}
					continue
				}

				o.Data = b[offset+1 : offset+int(o.Length)+1]
				objs = append(objs, &o)
				offset += int(o.Length) + 1

				if offset >= len(b) {
					break
				}
			}
			a.Objects = objs
		}
	case ComplexAck, SimpleAck, Error:
		a.InvokeID = b[offset]
		offset++
		a.Service = b[offset]
		offset++
		if len(b) > 3 {
			objs := []objects.APDUPayload{}
			for {
				o := objects.Object{
					TagNumber: b[offset] >> 4,
					TagClass:  common.IntToBool(int(b[offset]) & 0x8 >> 3),
					Length:    b[offset] & 0x7,
				}

				// Drop tags so that they don't get in the way!
				if b[offset] == objects.TagOpening || b[offset] == objects.TagClosing {
					offset++
					if offset >= len(b) {
						break
					}
					continue
				}

				o.Data = b[offset+1 : offset+int(o.Length)+1]
				objs = append(objs, &o)
				offset += int(o.Length) + 1

				if offset >= len(b) {
					break
				}
			}
			a.Objects = objs
		}
	}

	return nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (a *APDU) MarshalTo(b []byte) error {
	if len(b) < a.MarshalLen() {
		return common.ErrTooShortToMarshalBinary
	}

	var offset int = 0
	b[offset] = a.Type<<4 | a.Flags
	offset++

	switch a.Type {
	case UnConfirmedReq:
		b[offset] = a.Service
		offset++
		if a.MarshalLen() > 2 {
			for _, o := range a.Objects {
				ob, err := o.MarshalBinary()
				if err != nil {
					return err
				}

				copy(b[offset:offset+o.MarshalLen()], ob)
				offset += int(o.MarshalLen())

				if offset > a.MarshalLen() {
					return common.ErrTooShortToMarshalBinary
				}
			}
		}
	case ComplexAck, SimpleAck, Error:
		b[offset] = a.InvokeID
		offset++
		b[offset] = a.Service
		offset++
		if a.MarshalLen() > 4 {
			for _, o := range a.Objects {
				ob, err := o.MarshalBinary()
				if err != nil {
					return err
				}

				copy(b[offset:offset+o.MarshalLen()], ob)
				offset += o.MarshalLen()

				if offset > a.MarshalLen() {
					return common.ErrTooShortToMarshalBinary
				}
			}
		}
	case ConfirmedReq:
		b[offset] |= (a.MaxSeg & 0x7 << 4) | (a.MaxSize & 0xF)
		offset++
		b[offset] = a.InvokeID
		offset++
		b[offset] = a.Service
		offset++
		if a.MarshalLen() > 4 {
			for _, o := range a.Objects {
				ob, err := o.MarshalBinary()
				if err != nil {
					return err
				}

				copy(b[offset:offset+o.MarshalLen()], ob)
				offset += o.MarshalLen()

				if offset > a.MarshalLen() {
					return common.ErrTooShortToMarshalBinary
				}
			}
		}
	}

	return nil
}

// MarshalLen returns the serial length of APDU.
func (a *APDU) MarshalLen() int {
	var l int = 0
	switch a.Type {
	case ConfirmedReq:
		l += 4
	case ComplexAck, SimpleAck, Error:
		l += 3
	case UnConfirmedReq:
		l += 2
	}

	for _, o := range a.Objects {
		l += o.MarshalLen()
	}
	return l
}

// SetAPDUFlags sets APDU Flags to APDU.
func (a *APDU) SetAPDUFlags(sa, moreSegments, segmentedReq bool) {
	a.Flags = uint8(
		common.BoolToInt(sa)<<1 | common.BoolToInt(moreSegments)<<2 | common.BoolToInt(segmentedReq)<<3,
	)
}
