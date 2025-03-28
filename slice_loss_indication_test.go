// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Packet = (*SliceLossIndication)(nil) // assert is a Packet

func TestSliceLossIndicationUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      SliceLossIndication
		WantError error
	}{
		{
			Name: "valid",
			Data: []byte{
				// SliceLossIndication
				0x82, 0xcd, 0x0, 0x3,
				// sender=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// media=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// nack 0xAAAA, 0x5555
				0x55, 0x50, 0x00, 0x2C,
			},
			Want: SliceLossIndication{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
				SLI:        []SLIEntry{{0xaaa, 0, 0x2C}},
			},
		},
		{
			Name: "short report",
			Data: []byte{
				0x81, 0xcd, 0x0, 0x2,
				// ssrc=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// report ends early
			},
			WantError: errPacketTooShort,
		},
		{
			Name: "wrong type",
			Data: []byte{
				// v=2, p=0, count=1, SR, len=7
				0x81, 0xc8, 0x0, 0x7,
				// ssrc=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// ssrc=0xbc5e9a40
				0xbc, 0x5e, 0x9a, 0x40,
				// fracLost=0, totalLost=0
				0x0, 0x0, 0x0, 0x0,
				// lastSeq=0x46e1
				0x0, 0x0, 0x46, 0xe1,
				// jitter=273
				0x0, 0x0, 0x1, 0x11,
				// lsr=0x9f36432
				0x9, 0xf3, 0x64, 0x32,
				// delay=150137
				0x0, 0x2, 0x4a, 0x79,
			},
			WantError: errWrongType,
		},
		{
			Name:      "nil",
			Data:      nil,
			WantError: errPacketTooShort,
		},
	} {
		var sli SliceLossIndication
		err := sli.Unmarshal(test.Data)
		assert.ErrorIsf(t, err, test.WantError, "Unmarshal %q", test.Name)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.Want, sli, "Unmarshal %q rr", test.Name)
	}
}

func TestSliceLossIndicationRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Report    SliceLossIndication
		WantError error
	}{
		{
			Name: "valid",
			Report: SliceLossIndication{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
				SLI:        []SLIEntry{{1, 0xAA, 0x1F}, {1034, 0x05, 0x6}},
			},
		},
	} {
		data, err := test.Report.Marshal()
		assert.ErrorIsf(t, err, test.WantError, "Marshal %q", test.Name)
		if err != nil {
			continue
		}

		var decoded SliceLossIndication
		assert.NoErrorf(t, decoded.Unmarshal(data), "Unmarshal %q", test.Name)
		assert.Equalf(t, test.Report, decoded, "%q sli round trip mismatch", test.Name)
	}
}
