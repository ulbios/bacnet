package plumbing

// APDU type
const (
	ConfirmedReq uint8 = iota
	UnConfirmedReq
	SimpleAck
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
