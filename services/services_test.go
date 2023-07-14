// Copyright 2020 bacnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package services_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ulbios/bacnet"
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/plumbing"
	"github.com/ulbios/bacnet/services"
)

type serializeable interface {
	MarshalBinary() ([]byte, error)
	MarshalLen() int
}

type testCase struct {
	description string
	structured  serializeable
	serialized  []byte
}

func TestUnconfirmedWhoIs(t *testing.T) {
	t.Helper()
	var testcases = []testCase{
		{
			description: "Unconfirmed request WhoIs frame",
			structured: services.NewUnconfirmedWhoIs(
				plumbing.NewBVLC(plumbing.BVLCFuncBroadcast),
				plumbing.NewNPDU(false, false, false, false),
			),
			serialized: []byte{
				0x81, 0x0b, 0x00, 0x08, // BVLC
				0x01, 0x00, // NPDU
				0x10, 0x08, // APDU
			},
		},
	}

	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				msg, err := bacnet.Parse(c.serialized)
				if err != nil {
					t.Fatal(err)
				}

				want, got := c.structured, msg
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			})
			t.Run("Serialize", func(t *testing.T) {
				b, err := c.structured.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				want, got := c.serialized, b
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			})
		})
	}
}

func TestUnconfirmedIAm(t *testing.T) {
	t.Helper()
	var testcases = []testCase{
		{
			description: "Unconfirmed request IAm frame",
			structured: services.NewUnconfirmedIAm(
				plumbing.NewBVLC(plumbing.BVLCFuncBroadcast),
				plumbing.NewNPDU(false, false, false, false),
			),
			serialized: []byte{
				0x81, 0x0b, 0x00, 0x14, // BVLC
				0x01, 0x00, // NPDU
				0x10, 0x00, // APDU
				0xc4, 0x02, 0x00, 0x00, 0x01, // device object
				0x22, 0x04, 0x00, // Max APDU length accepted
				0x91, 0x00, // Segmentation supported
				0x21, 0x01, // Vendor ID
			},
		},
	}

	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				msg, err := bacnet.Parse(c.serialized)
				if err != nil {
					t.Fatal(err)
				}

				want, got := c.structured, msg
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			})
			t.Run("Serialize", func(t *testing.T) {
				b, err := c.structured.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}

				want, got := c.serialized, b
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			})
		})
	}
}
func TestBoolToInt(t *testing.T) {
	cases := []struct {
		description string
		b           bool
		i           int
	}{
		{
			"case of true",
			true,
			1,
		},
		{
			"case of false",
			false,
			0,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			if got, want := common.BoolToInt(c.b), c.i; got != want {
				t.Fail()
			}

		})
	}
}

func TestIntToBool(t *testing.T) {
	cases := []struct {
		description string
		b           bool
		i           int
	}{
		{
			"case of true",
			true,
			1,
		},
		{
			"case of false",
			false,
			0,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			if got, want := common.IntToBool(c.i), c.b; got != want {
				t.Fail()
			}

		})
	}
}
