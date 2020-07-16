// Copyright 2020 bacnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package bacnet

import (
	"encoding/binary"
)

// BoolToInt converts bool to int.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IntToBool converts int to bool.
func IntToBool(n int) bool {
	return n != 0
}

// Message is an interface that defines BACnet messages.
type Message interface {
	MarshalBinary() ([]byte, error)
	MarshalTo([]byte) error
	UnmarshalBinary([]byte) error
	MarshalLen() int
}

// BVLCType is used for BACnet/IP in BVLL.
const BVLCType = 0x81

// BVLCFunc determines unicast or broadcast in BVLL.
const (
	BVLCFuncUnicast   = 0x0a
	BVLCFuncBroadcast = 0x0b
)

// BVLC is a BVLC frame.
type BVLC struct {
	Type     uint8
	Function uint8
	Length   uint16
}

// NewBVLC creates a BVLC.
func NewBVLC(f uint8) *BVLC {
	bvlc := &BVLC{
		Type:     BVLCType,
		Function: f,
	}
	return bvlc
}

// UnmarshalBinary sets the values retrieved from byte sequence in a BVLC frame.
func (bvlc *BVLC) UnmarshalBinary(b []byte) error {
	if l := len(b); l < bvlc.MarshalLen() {
		return ErrTooShortToParse
	}
	bvlc.Type = b[0]
	bvlc.Function = b[1]
	bvlc.Length = binary.BigEndian.Uint16(b[2:4])

	return nil
}

// MarshalBinary returns the byte sequence generated from a BVLC instance.
func (bvlc *BVLC) MarshalBinary() ([]byte, error) {
	b := make([]byte, bvlc.MarshalLen())
	if err := bvlc.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

const bvlclen = 4

// MarshalLen returns the serial length of BVLC.
func (bvlc *BVLC) MarshalLen() int {
	return bvlclen
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (bvlc *BVLC) MarshalTo(b []byte) error {
	if len(b) < bvlc.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}
	b[0] = byte(bvlc.Type)
	b[1] = byte(bvlc.Function)
	binary.BigEndian.PutUint16(b[2:4], bvlc.Length)
	return nil
}

// SetControlFlags sets control flags to NPDU.
func (n *NPDU) SetControlFlags(nsduContain bool, dstSpecifier bool, srcSpecifier bool, expectingReply bool) {
	n.Control = uint8(
		BoolToInt(nsduContain)<<7 | BoolToInt(dstSpecifier)<<5 | BoolToInt(srcSpecifier)<<3 | BoolToInt(expectingReply)<<2,
	)
}

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

// UnmarshalBinary sets the values retrieved from byte sequence in a NPDU frame.
func (n *NPDU) UnmarshalBinary(b []byte) error {
	if l := len(b); l < n.MarshalLen() {
		return ErrTooShortToParse
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
		return ErrTooShortToMarshalBinary
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

// APDU is a Application protocol DAta Units.
type APDU struct {
	Type     uint8
	Flags    uint8
	MaxSeg   uint8
	MaxSize  uint8
	InvokeID uint8
	Service  uint8
	Objects  []Object
}

// NewAPDU creates an APDU.
func NewAPDU(t, c uint8, o []Object) *APDU {
	return &APDU{
		Type:    t,
		Service: c,
		Objects: o,
	}
}

// UnmarshalBinary sets the values retrieved from byte sequence in a APDU frame.
func (a *APDU) UnmarshalBinary(b []byte) error {
	if l := len(b); l < a.MarshalLen() {
		return ErrTooShortToParse
	}

	a.Type = b[0] >> 4
	a.Flags = b[0] & 0x7

	var offset int = 1

	switch a.Type {
	case UnConfirmedReq:
		b[offset] = a.Service
		offset++
		if len(b) > 2 {
			objs := []Object{}
			for {
				o := Object{
					TagNumber: b[offset] >> 4,
					TagClass:  IntToBool(int(b[offset]) & 0x8 >> 3),
					Length:    b[offset+1],
				}
				o.Data = b[offset+2 : o.Length]
				objs = append(objs, o)
				offset++

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
		return ErrTooShortToMarshalBinary
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
				copy(b[offset:o.MarshalLen()], ob)
				offset++

				if offset >= a.MarshalLen() {
					return ErrTooShortToMarshalBinary
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
	case UnConfirmedReq:
		l += 2
	}

	for _, o := range a.Objects {
		l += int(1 + o.Length)
	}
	return l
}

// APDU type
const (
	ConfirmedReq uint8 = iota
	UnConfirmedReq
	ComplexAck
	SegmentAck
	Error
	Reject
	Abort
)

// APDU flags for confirmedRequest
const (
	SA uint8 = (iota + 1) * 2
	MoreSegments
	SegmentedRequest
)

// SetAPDUFlags sets APDU Flags to APDU.
func (a *APDU) SetAPDUFlags(sa, moreSegments, segmentedReq bool) {
	a.Flags = uint8(
		BoolToInt(sa)<<1 | BoolToInt(moreSegments)<<2 | BoolToInt(segmentedReq)<<3,
	)
}

// Object is an object in APDU.
type Object struct {
	TagNumber uint8
	TagClass  bool
	Length    uint8
	Data      []byte
}

// NewObject creates an Object.
func NewObject(n uint8, c bool, d []byte) *Object {
	return &Object{
		TagNumber: n,
		TagClass:  c,
		Length:    uint8(len(d)),
		Data:      d,
	}
}

const objLenMin int = 2

// UnmarshalBinary sets the values retrieved from byte sequence in a Object frame.
func (o *Object) UnmarshalBinary(b []byte) error {
	if l := len(b); l < objLenMin {
		return ErrTooShortToParse
	}
	o.TagNumber = b[0] >> 4
	o.TagClass = IntToBool(int(b[0]) & 0x8 >> 3)
	o.Length = b[0] & 0x7

	if l := len(b); l < int(o.Length) {
		return ErrTooShortToParse
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
		return ErrTooShortToMarshalBinary
	}
	b[0] = o.TagNumber<<4 | uint8(BoolToInt(o.TagClass))<<3 | o.Length
	if o.Length > 0 {
		copy(b[1:o.Length+1], o.Data)
	}
	return nil
}

// MarshalLen returns the serial length of Object.
func (o *Object) MarshalLen() int {
	return 1 + int(o.Length)
}

// UnconfirmedWhoIs is a BACnet message.
type UnconfirmedWhoIs struct {
	*BVLC
	*NPDU
	*APDU
	*Object
}

// NewUnconfirmedWhoIs creates a UnconfirmedWhoIs.
func NewUnconfirmedWhoIs(bvlc *BVLC, npdu *NPDU) *UnconfirmedWhoIs {
	u := &UnconfirmedWhoIs{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: NewAPDU(UnConfirmedReq, ServiceUnconfirmedWhoIs, nil),
	}
	u.SetLength()
	return u
}

// UnmarshalBinary sets the values retrieved from byte sequence in a UnconfirmedWhoIs frame.
func (u *UnconfirmedWhoIs) UnmarshalBinary(b []byte) error {
	if l := len(b); l < u.MarshalLen() {
		return ErrTooShortToParse
	}

	var offset int = 0
	if err := u.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return ErrTooShortToParse
	}
	offset += u.BVLC.MarshalLen()

	if err := u.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return ErrTooShortToParse
	}
	offset += u.NPDU.MarshalLen()

	if err := u.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return ErrTooShortToParse
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
		return ErrTooShortToMarshalBinary
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

// MarshalLen returns the serial length of WhoIs.
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

// BACnet is an interface defines BACnet messages.
type BACnet interface {
	MarshalBinary() ([]byte, error)
	MarshalTo([]byte) error
	UnmarshalBinary([]byte) error
	MarshalLen() int
}

// Services in APDU of which type is unconfirmed request.
const (
	ServiceUnconfirmedIAm uint8 = iota
	ServiceUnconfirmedIHave
	ServiceUnconfirmedCOVNotification
	ServiceUnconfirmedEventNotification
	ServiceUnconfirmedPrivateTransfer
	ServiceUnconfirmedTextMessage
	ServiceUnconfirmedTimeSync
	ServiceUnconfirmedWhoHas
	ServiceUnconfirmedWhoIs
	ServiceUnconfirmedUTCTimeSync
	ServiceUnconfirmedWriteGroup
)

// Services in APDU of which type is confirmed request.
const (
	ServiceConfirmedAcknowledgeAlarm uint8 = iota
	ServiceConfirmedCOVNotification
	ServiceConfirmedEventNotification
	ServiceConfirmedGetAlarmSummary
	ServiceConfirmedGetEnrollmentSummary
	ServiceConfirmedSubscribeCOV
	ServiceConfirmedAtomicReadFile
	ServiceConfirmedAtomicWriteFile
	ServiceConfirmedAddListElement
	ServiceConfirmedRemoveListElement
	ServiceConfirmedCreateObject
	ServiceConfirmedDeleteObject
	ServiceConfirmedReadProperty
	ServiceConfirmedReadPropConditional
	ServiceConfirmedReadPropMultiple
	ServiceConfirmedWriteProperty
	ServiceConfirmedWritePropMultiple
	ServiceConfirmedDeviceCommunicationControl
	ServiceConfirmedPrivateTransfer
	ServiceConfirmedTextMessage
	ServiceConfirmedReinitializeDevice
	ServiceConfirmedVTOpen
	ServiceConfirmedVTClose
	ServiceConfirmedVTData
	ServiceConfirmedAuthenticate
	ServiceConfirmedRequestKey
)

const bacnetLenMin = 8

// Parse decodes the given bytes.
func Parse(b []byte) (BACnet, error) {
	if len(b) < bacnetLenMin {
		return nil, ErrTooShortToParse
	}

	var bacnet BACnet

	combine := func(t, s uint8) uint16 {
		return uint16(t)<<8 | uint16(s)
	}

	var offset = 0
	bvlc := NewBVLC(BVLCFuncBroadcast)
	offset += bvlc.MarshalLen()

	npdu := NewNPDU(false, false, false, false)
	offset += npdu.MarshalLen()

	c := combine(b[offset], b[offset+1])

	switch c {
	case combine(UnConfirmedReq<<4, ServiceUnconfirmedWhoIs):
		bacnet = NewUnconfirmedWhoIs(bvlc, npdu)
	case combine(UnConfirmedReq<<4, ServiceUnconfirmedIAm):
		// TODO: bacnet = NewUnconfirmedIam(bvlc, npdu)
	default:
		return nil, ErrNotImplemented
	}

	if err := bacnet.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return bacnet, nil
}
