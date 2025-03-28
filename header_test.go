// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      Header
		WantError error
	}{
		{
			Name: "valid",
			Data: []byte{
				// v=2, p=0, count=1, RR, len=7
				0x81, 0xc9, 0x00, 0x07,
			},
			Want: Header{
				Padding: false,
				Count:   1,
				Type:    TypeReceiverReport,
				Length:  7,
			},
		},
		{
			Name: "also valid",
			Data: []byte{
				// v=2, p=1, count=1, BYE, len=7
				0xa1, 0xcc, 0x00, 0x07,
			},
			Want: Header{
				Padding: true,
				Count:   1,
				Type:    TypeApplicationDefined,
				Length:  7,
			},
		},
		{
			Name: "bad version",
			Data: []byte{
				// v=0, p=0, count=0, RR, len=4
				0x00, 0xc9, 0x00, 0x04,
			},
			WantError: errBadVersion,
		},
	} {
		var h Header
		err := h.Unmarshal(test.Data)
		assert.ErrorIsf(t, err, test.WantError, "Unmarshal %q header mispmatch", test.Name)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.Want, h, "Unmarshal %q header mismatch", test.Name)
	}
}

func TestHeaderRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Header    Header
		WantError error
	}{
		{
			Name: "valid",
			Header: Header{
				Padding: true,
				Count:   31,
				Type:    TypeSenderReport,
				Length:  4,
			},
		},
		{
			Name: "also valid",
			Header: Header{
				Padding: false,
				Count:   28,
				Type:    TypeReceiverReport,
				Length:  65535,
			},
		},
		{
			Name: "invalid count",
			Header: Header{
				Count: 40,
			},
			WantError: errInvalidHeader,
		},
	} {
		data, err := test.Header.Marshal()
		assert.ErrorIsf(t, err, test.WantError, "Marshal %q", test.Name)
		if err != nil {
			continue
		}

		var decoded Header
		assert.NoErrorf(t, decoded.Unmarshal(data), "Unmarshal %q", test.Name)
		assert.Equalf(t, test.Header, decoded, "%q header round trip mismatch", test.Name)
	}
}
