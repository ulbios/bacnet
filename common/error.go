// Copyright 2020 bacnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package common

import "errors"

// Error definitions.
var (
	ErrTooShortToMarshalBinary = errors.New("insufficient buffer to serialize parameter to")
	ErrTooShortToParse         = errors.New("too short to decode as parameter")
	ErrNotImplemented          = errors.New("not implemented type")
	ErrTooBigValue             = errors.New("too big value")
	ErrWrongTagNumber          = errors.New("wrong tag number")
	ErrWrongObjectCount        = errors.New("wrong object count")
	ErrWrongStructure          = errors.New("unexpected object structure")
	ErrWrongPayload            = errors.New("wrong payload type")
)
