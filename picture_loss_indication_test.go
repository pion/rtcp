// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Packet = (*PictureLossIndication)(nil) // assert is a Packet

func TestPictureLossIndicationUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      PictureLossIndication
		WantError error
	}{
		{
			Name: "valid",
			Data: []byte{
				// v=2, p=0, FMT=1, PSFB, len=1
				0x81, 0xce, 0x00, 0x02,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
			},
			Want: PictureLossIndication{
				SenderSSRC: 0x0,
				MediaSSRC:  0x4bc4fcb4,
			},
		},
		{
			Name: "packet too short",
			Data: []byte{
				0x00, 0x00, 0x00, 0x00,
			},
			WantError: errPacketTooShort,
		},
		{
			Name: "invalid header",
			Data: []byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			WantError: errBadVersion,
		},
		{
			Name: "wrong type",
			Data: []byte{
				// v=2, p=0, FMT=1, RR, len=1
				0x81, 0xc9, 0x00, 0x02,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
			},
			WantError: errWrongType,
		},
		{
			Name: "wrong fmt",
			Data: []byte{
				// v=2, p=0, FMT=2, RR, len=1
				0x82, 0xc9, 0x00, 0x02,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
			},
			WantError: errWrongType,
		},
	} {
		var pli PictureLossIndication
		err := pli.Unmarshal(test.Data)
		assert.ErrorIsf(t, err, test.WantError, "Unmarshal %q", test.Name)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.Want, pli, "Unmarshal %q", test.Name)
	}
}

func TestPictureLossIndicationRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Packet    PictureLossIndication
		WantError error
	}{
		{
			Name: "valid",
			Packet: PictureLossIndication{
				SenderSSRC: 1,
				MediaSSRC:  2,
			},
		},
		{
			Name: "also valid",
			Packet: PictureLossIndication{
				SenderSSRC: 5000,
				MediaSSRC:  6000,
			},
		},
	} {
		data, err := test.Packet.Marshal()
		assert.ErrorIsf(t, err, test.WantError, "Marshal %q", test.Name)
		if err != nil {
			continue
		}

		var decoded PictureLossIndication
		assert.NoErrorf(t, decoded.Unmarshal(data), "Unmarshal %q", test.Name)
		assert.Equalf(t, test.Packet, decoded, "%q rr round trip mismatch", test.Name)
	}
}

func TestPictureLossIndicationUnmarshalHeader(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      Header
		WantError error
	}{
		{
			Name: "valid header",
			Data: []byte{
				// v=2, p=0, FMT=1, PSFB, len=1
				0x81, 0xce, 0x00, 0x02,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
			},
			Want: Header{
				Count:  FormatPLI,
				Type:   TypePayloadSpecificFeedback,
				Length: pliLength,
			},
		},
	} {
		var pli PictureLossIndication
		err := pli.Unmarshal(test.Data)
		assert.ErrorIsf(t, err, test.WantError, "Unmarshal header %q", test.Name)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.Want, pli.Header(), "Unmarshal header %q", test.Name)
	}
}
