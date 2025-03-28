// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Packet = (*TransportLayerNack)(nil) // assert is a Packet

func TestTransportLayerNackUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      TransportLayerNack
		WantError error
	}{
		{
			Name: "valid",
			Data: []byte{
				// TransportLayerNack
				0x81, 0xcd, 0x0, 0x3,
				// sender=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// media=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// nack 0xAAAA, 0x5555
				0xaa, 0xaa, 0x55, 0x55,
			},
			Want: TransportLayerNack{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
				Nacks:      []NackPair{{0xaaaa, 0x5555}},
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
			Name: "bad length",
			Data: []byte{
				// TransportLayerNack
				0x81, 0xcd, 0x0, 0x2,
				// sender=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// media=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
			},
			WantError: errBadLength,
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
		var tln TransportLayerNack
		err := tln.Unmarshal(test.Data)
		assert.ErrorIsf(t, err, test.WantError, "Unmarshal %q", test.Name)
		assert.Equalf(t, test.Want, tln, "Unmarshal %q", test.Name)
	}
}

func TestTransportLayerNackRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Report    TransportLayerNack
		WantError error
	}{
		{
			Name: "valid",
			Report: TransportLayerNack{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
				Nacks:      []NackPair{{1, 0xAA}, {1034, 0x05}},
			},
		},
	} {
		data, err := test.Report.Marshal()
		assert.ErrorIsf(t, err, test.WantError, "Marshal err on %q", test.Name)

		var decoded TransportLayerNack
		assert.NoErrorf(t, decoded.Unmarshal(data), "Unmarshal %q", test.Name)
		assert.Equalf(t, test.Report, decoded, "Unmarshal %q: decoded mismatch", test.Name)
	}
}

func testNackPair(t *testing.T, s []uint16, n NackPair) {
	t.Helper()

	l := n.PacketList()
	assert.Equalf(t, s, l, "NackPair %v mismatch", n)
}

func TestNackPair(t *testing.T) {
	testNackPair(t, []uint16{42}, NackPair{42, 0})
	testNackPair(t, []uint16{42, 43}, NackPair{42, 1})
	testNackPair(t, []uint16{42, 44}, NackPair{42, 2})
	testNackPair(t, []uint16{42, 43, 44}, NackPair{42, 3})
	testNackPair(t, []uint16{42, 42 + 16}, NackPair{42, 0x8000})
}

func TestNackPairRange(t *testing.T) {
	pair := NackPair{42, 2}

	out := make([]uint16, 0)
	pair.Range(func(s uint16) bool {
		out = append(out, s)

		return true
	})
	assert.Equal(t, []uint16{42, 44}, out)

	out = make([]uint16, 0)
	pair.Range(func(s uint16) bool {
		out = append(out, s)

		return false
	})
	assert.Equal(t, []uint16{42}, out)
}

func TestTransportLayerNackPairGeneration(t *testing.T) {
	for _, test := range []struct {
		Name            string
		SequenceNumbers []uint16
		Expected        []NackPair
	}{
		{
			"No Sequence Numbers",
			[]uint16{},
			[]NackPair{},
		},
		{
			"Single Sequence Number",
			[]uint16{100},
			[]NackPair{
				{PacketID: 100, LostPackets: 0x0},
			},
		},
		{
			"Multiple in range, Single NACKPair",
			[]uint16{100, 101, 105, 115},
			[]NackPair{
				{PacketID: 100, LostPackets: 0x4011},
			},
		},
		{
			"Multiple Ranges, Multiple NACKPair",
			[]uint16{100, 117, 500, 501, 502},
			[]NackPair{
				{PacketID: 100, LostPackets: 0},
				{PacketID: 117, LostPackets: 0},
				{PacketID: 500, LostPackets: 0x3},
			},
		},
	} {
		actual := NackPairsFromSequenceNumbers(test.SequenceNumbers)
		assert.Equalf(t, test.Expected, actual, "%q NackPair generation mismatch", test.Name)
	}
}
