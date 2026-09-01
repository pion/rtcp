// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// An RTCP packet from a packet dump.
func realPacket() []byte {
	return []byte{
		// Receiver Report (offset=0)
		// v=2, p=0, count=1, RR, len=7
		0x81, 0xc9, 0x0, 0x7,
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

		// Source Description (offset=32)
		// v=2, p=0, count=1, SDES, len=12
		0x81, 0xca, 0x0, 0xc,
		// ssrc=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,
		// CNAME, len=38
		0x1, 0x26,
		// text="{9c00eb92-1afb-9d49-a47d-91f64eee69f5}"
		0x7b, 0x39, 0x63, 0x30,
		0x30, 0x65, 0x62, 0x39,
		0x32, 0x2d, 0x31, 0x61,
		0x66, 0x62, 0x2d, 0x39,
		0x64, 0x34, 0x39, 0x2d,
		0x61, 0x34, 0x37, 0x64,
		0x2d, 0x39, 0x31, 0x66,
		0x36, 0x34, 0x65, 0x65,
		0x65, 0x36, 0x39, 0x66,
		0x35, 0x7d,
		// END + padding
		0x0, 0x0, 0x0, 0x0,

		// Goodbye (offset=84)
		// v=2, p=0, count=1, BYE, len=1
		0x81, 0xcb, 0x0, 0x1,
		// source=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,

		// Picture Loss Indication (offset=92)
		0x81, 0xce, 0x0, 0x2,
		// sender=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,
		// media=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,

		// RapidResynchronizationRequest (offset=104)
		0x85, 0xcd, 0x0, 0x2,
		// sender=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,
		// media=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,

		// ApplicationDefined (offset=116)
		0x80, 0xcc, 0x00, 0x03,
		// sender=0x4baae1ab
		0x4b, 0xaa, 0xe1, 0xab,
		// name='NAME'
		0x4E, 0x41, 0x4D, 0x45,
		// data='ABCD'
		0x41, 0x42, 0x43, 0x44,
	}
}

func TestUnmarshal(t *testing.T) {
	packet, err := Unmarshal(realPacket())
	assert.NoError(t, err)

	expected := []Packet{
		&ReceiverReport{
			SSRC: 0x902f9e2e,
			Reports: []ReceptionReport{{
				SSRC:               0xbc5e9a40,
				FractionLost:       0,
				TotalLost:          0,
				LastSequenceNumber: 0x46e1,
				Jitter:             273,
				LastSenderReport:   0x9f36432,
				Delay:              150137,
			}},
			ProfileExtensions: []byte{},
		},
		NewCNAMESourceDescription(0x902f9e2e, "{9c00eb92-1afb-9d49-a47d-91f64eee69f5}"),
		&Goodbye{
			Sources: []uint32{0x902f9e2e},
		},
		&PictureLossIndication{
			SenderSSRC: 0x902f9e2e,
			MediaSSRC:  0x902f9e2e,
		},
		&RapidResynchronizationRequest{
			SenderSSRC: 0x902f9e2e,
			MediaSSRC:  0x902f9e2e,
		},
		&ApplicationDefined{
			SSRC: 0x4baae1ab,
			Name: "NAME",
			Data: []byte{0x41, 0x42, 0x43, 0x44},
		},
	}

	assert.Equal(t, expected, packet)
}

func TestUnmarshalNil(t *testing.T) {
	_, err := Unmarshal(nil)
	assert.ErrorIs(t, err, errInvalidHeader)
}

func TestInvalidHeaderLength(t *testing.T) {
	invalidPacket := []byte{
		// Receiver Report (offset=0)
		// v=2, p=0, count=1, RR, len=100
		0x81, 0xc9, 0x0, 0x64,
	}

	_, err := Unmarshal(invalidPacket)
	assert.ErrorIs(t, err, errPacketTooShort)
}

func TestUnmarshalMaxLengthRawPacket(t *testing.T) {
	rawPacket := make([]byte, 4*(0xFFFF+1))
	rawPacket[0] = 0x80
	rawPacket[1] = 0xfa
	rawPacket[2] = 0xff
	rawPacket[3] = 0xff

	packets, err := Unmarshal(rawPacket)
	assert.NoError(t, err)
	assert.Len(t, packets, 1)

	parsed, ok := packets[0].(*RawPacket)
	assert.True(t, ok)
	assert.Equal(t, rawPacket, []byte(*parsed))
}

// Sender and receiver reports carry profile-specific extensions that are padded
// out to a whole 32-bit word, so both the header length field and MarshalSize
// have to account for that padding.
func TestReportProfileExtensionsRoundTrip(t *testing.T) {
	for _, extLength := range []int{1, 4, 5, 8, 24} {
		extensions := bytes.Repeat([]byte{0xab}, extLength)
		padded := append(bytes.Repeat([]byte{0xab}, extLength), make([]byte, getPadding(extLength))...)
		reports := []ReceptionReport{{SSRC: 2}}

		receiverReport := &ReceiverReport{SSRC: 1, Reports: reports, ProfileExtensions: extensions}
		senderReport := &SenderReport{SSRC: 1, Reports: reports, ProfileExtensions: extensions}

		for _, test := range []struct {
			Name   string
			Packet Packet
			Header func() Header
			Want   Packet
		}{
			{
				Name:   "ReceiverReport",
				Packet: receiverReport,
				Header: receiverReport.Header,
				Want:   &ReceiverReport{SSRC: 1, Reports: reports, ProfileExtensions: padded},
			},
			{
				Name:   "SenderReport",
				Packet: senderReport,
				Header: senderReport.Header,
				Want:   &SenderReport{SSRC: 1, Reports: reports, ProfileExtensions: padded},
			},
		} {
			data, err := test.Packet.Marshal()
			assert.NoErrorf(t, err, "Marshal %s with %d extension bytes", test.Name, extLength)

			assert.Lenf(t, data, test.Packet.MarshalSize(),
				"MarshalSize %s with %d extension bytes", test.Name, extLength)
			assert.Equalf(t, uint16(len(data)/4-1), test.Header().Length, //nolint:gosec // G115
				"header length %s with %d extension bytes", test.Name, extLength)

			packets, err := Unmarshal(data)
			assert.NoErrorf(t, err, "Unmarshal %s with %d extension bytes", test.Name, extLength)
			if !assert.Lenf(t, packets, 1, "packet count %s with %d extension bytes", test.Name, extLength) {
				continue
			}
			assert.Equalf(t, test.Want, packets[0],
				"round trip %s with %d extension bytes", test.Name, extLength)
		}
	}
}
