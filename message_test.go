// Copyright 2020 bacnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package bacnet_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kazukiigeta/bacnet"
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
			structured: bacnet.NewUnconfirmedWhoIs(
				bacnet.NewBVLC(bacnet.BVLCFuncBroadcast),
				bacnet.NewNPDU(false, false, false, false),
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

				got, want := msg, c.structured
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("differs: (+got -want)\n%s", diff)
				}
			})
			t.Run("Serialize", func(t *testing.T) {
				b, err := c.structured.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				got, want := b, c.serialized
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("differs: (+got -want)\n%s", diff)
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
			if got, want := bacnet.BoolToInt(c.b), c.i; got != want {
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
			if got, want := bacnet.IntToBool(c.i), c.b; got != want {
				t.Fail()
			}

		})
	}
}
