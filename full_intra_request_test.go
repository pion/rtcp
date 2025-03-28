// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullIntraRequestUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      FullIntraRequest
		WantError error
	}{
		{
			Name: "valid",
			Data: []byte{
				// v=2, p=0, FMT=4, PSFB, len=4
				0x84, 0xce, 0x00, 0x04,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x12345678
				0x12, 0x34, 0x56, 0x78,
				// Seqno=0x42
				0x42, 0x00, 0x00, 0x00,
			},
			Want: FullIntraRequest{
				SenderSSRC: 0x0,
				MediaSSRC:  0x4bc4fcb4,
				FIR: []FIREntry{
					{
						SSRC:           0x12345678,
						SequenceNumber: 0x42,
					},
				},
			},
		},
		{
			Name: "also valid",
			Data: []byte{
				// v=2, p=0, FMT=4, PSFB, len=6
				0x84, 0xce, 0x00, 0x06,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x12345678
				0x12, 0x34, 0x56, 0x78,
				// Seqno=0x42
				0x42, 0x00, 0x00, 0x00,
				// ssrc=0x98765432
				0x98, 0x76, 0x54, 0x32,
				// Seqno=0x57
				0x57, 0x00, 0x00, 0x00,
			},
			Want: FullIntraRequest{
				SenderSSRC: 0x0,
				MediaSSRC:  0x4bc4fcb4,
				FIR: []FIREntry{
					{
						SSRC:           0x12345678,
						SequenceNumber: 0x42,
					},
					{
						SSRC:           0x98765432,
						SequenceNumber: 0x57,
					},
				},
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
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			WantError: errBadVersion,
		},
		{
			Name: "wrong type",
			Data: []byte{
				// v=2, p=0, FMT=4, RR, len=4
				0x84, 0xc9, 0x00, 0x04,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x12345678
				0x12, 0x34, 0x56, 0x78,
				// Seqno=0x42
				0x42, 0x00, 0x00, 0x00,
			},
			WantError: errWrongType,
		},
		{
			Name: "wrong fmt",
			Data: []byte{
				// v=2, p=0, FMT=2, PSFB, len=4
				0x82, 0xce, 0x00, 0x04,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x12345678
				0x12, 0x34, 0x56, 0x78,
				// Seqno=0x42
				0x42, 0x00, 0x00, 0x00,
			},
			WantError: errWrongType,
		},
		{
			Name: "wrong length",
			Data: []byte{
				// v=2, p=0, FMT=4, PSFB, len=3
				0x84, 0xce, 0x00, 0x03,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x12345678
				0x12, 0x34, 0x56, 0x78,
			},
			WantError: errBadLength,
		},
	} {
		var fir FullIntraRequest
		err := fir.Unmarshal(test.Data)
		assert.ErrorIsf(t, err, test.WantError, "Unmarshal %q rr mismatch", test.Name)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.Want, fir, "Unmarshal %q rr mismatch", test.Name)
	}
}

func TestFullIntraRequestRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Packet    FullIntraRequest
		WantError error
	}{
		{
			Name: "valid",
			Packet: FullIntraRequest{
				SenderSSRC: 1,
				MediaSSRC:  2,
				FIR: []FIREntry{{
					SSRC:           3,
					SequenceNumber: 42,
				}},
			},
		},
		{
			Name: "also valid",
			Packet: FullIntraRequest{
				SenderSSRC: 5000,
				MediaSSRC:  6000,
				FIR: []FIREntry{{
					SSRC:           3,
					SequenceNumber: 57,
				}},
			},
		},
	} {
		data, err := test.Packet.Marshal()
		assert.ErrorIsf(t, err, test.WantError, "Marshal %q", test.Name)
		if err != nil {
			continue
		}

		var decoded FullIntraRequest
		assert.NoErrorf(t, decoded.Unmarshal(data), "Unmarshal %q", test.Name)
		assert.Equalf(t, test.Packet, decoded, "%q rr header mismatch", test.Name)
	}
}

func TestFullIntraRequestUnmarshalHeader(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      Header
		WantError error
	}{
		{
			Name: "valid header",
			Data: []byte{
				// v=2, p=0, FMT=1, PSFB, len=4
				0x84, 0xce, 0x00, 0x04,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x00000000
				0x00, 0x00, 0x00, 0x00,
				// Seqno=0x22
				0x22, 0x00, 0x00, 0x00,
			},
			Want: Header{
				Count:  FormatFIR,
				Type:   TypePayloadSpecificFeedback,
				Length: 4,
			},
		},
	} {
		var fir FullIntraRequest
		err := fir.Unmarshal(test.Data)
		assert.ErrorIsf(t, err, test.WantError, "Unmarshal header %q rr mismatch", test.Name)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.Want, fir.Header(), "Unmarshal header %q rr mismatch", test.Name)
	}
}
