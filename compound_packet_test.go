// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Packet = (*CompoundPacket)(nil) // assert is a Packet

func TestReadEOF(t *testing.T) {
	shortHeader := []byte{
		0x81, 0xc9, // missing type & len
	}

	_, err := Unmarshal(shortHeader)
	assert.Error(t, err)
}

func TestBadCompound(t *testing.T) {
	// trailing data!
	badcompound := realPacket()[:34]
	packets, err := Unmarshal(badcompound)
	assert.Error(t, err)

	assert.Nil(t, packets)

	badcompound = realPacket()[84:104]

	packets, err = Unmarshal(badcompound)
	assert.NoError(t, err)

	compound := CompoundPacket(packets)

	// this should return an error,
	// it violates the "must start with RR or SR" rule
	assert.ErrorIs(t, compound.Validate(), errBadFirstPacket)
	assert.Equal(t, 2, len(compound))

	_, ok := compound[0].(*Goodbye)
	assert.True(t, ok)

	_, ok = compound[1].(*PictureLossIndication)
	assert.True(t, ok)
}

func TestValidPacket(t *testing.T) {
	cname := NewCNAMESourceDescription(1234, "cname")

	for _, test := range []struct {
		Name   string
		Packet CompoundPacket
		Err    error
	}{
		{
			Name:   "empty",
			Packet: CompoundPacket{},
			Err:    errEmptyCompound,
		},
		{
			Name: "no cname",
			Packet: CompoundPacket{
				&SenderReport{},
			},
			Err: errMissingCNAME,
		},
		{
			Name: "just BYE",
			Packet: CompoundPacket{
				&Goodbye{},
			},
			Err: errBadFirstPacket,
		},
		{
			Name: "SDES / no cname",
			Packet: CompoundPacket{
				&SenderReport{},
				&SourceDescription{},
			},
			Err: errMissingCNAME,
		},
		{
			Name: "just SR",
			Packet: CompoundPacket{
				&SenderReport{},
				cname,
			},
			Err: nil,
		},
		{
			Name: "multiple SRs",
			Packet: CompoundPacket{
				&SenderReport{},
				&SenderReport{},
				cname,
			},
			Err: errPacketBeforeCNAME,
		},
		{
			Name: "just RR",
			Packet: CompoundPacket{
				&ReceiverReport{},
				cname,
			},
			Err: nil,
		},
		{
			Name: "multiple RRs",
			Packet: CompoundPacket{
				&ReceiverReport{},
				&ReceiverReport{},
				cname,
			},
			Err: nil,
		},
		{
			Name: "goodbye",
			Packet: CompoundPacket{
				&ReceiverReport{},
				cname,
				&Goodbye{},
			},
			Err: nil,
		},
	} {
		assert.ErrorIsf(t, test.Packet.Validate(), test.Err, "Validate(%s)", test.Name)
	}
}

func TestCNAME(t *testing.T) {
	cname := NewCNAMESourceDescription(1234, "cname")

	for _, test := range []struct {
		Name   string
		Packet CompoundPacket
		Err    error
		Text   string
	}{
		{
			Name: "no cname",
			Packet: CompoundPacket{
				&SenderReport{},
			},
			Err: errMissingCNAME,
		},
		{
			Name: "SDES / no cname",
			Packet: CompoundPacket{
				&SenderReport{},
				&SourceDescription{},
			},
			Err: errMissingCNAME,
		},
		{
			Name: "just SR",
			Packet: CompoundPacket{
				&SenderReport{},
				cname,
			},
			Err:  nil,
			Text: "cname",
		},
		{
			Name: "multiple SRs",
			Packet: CompoundPacket{
				&SenderReport{},
				&SenderReport{},
				cname,
			},
			Err:  errPacketBeforeCNAME,
			Text: "cname",
		},
		{
			Name: "just RR",
			Packet: CompoundPacket{
				&ReceiverReport{},
				cname,
			},
			Err:  nil,
			Text: "cname",
		},
		{
			Name: "multiple RRs",
			Packet: CompoundPacket{
				&ReceiverReport{},
				&ReceiverReport{},
				cname,
			},
			Err:  nil,
			Text: "cname",
		},
		{
			Name: "goodbye",
			Packet: CompoundPacket{
				&ReceiverReport{},
				cname,
				&Goodbye{},
			},
			Err:  nil,
			Text: "cname",
		},
	} {
		assert.ErrorIsf(t, test.Packet.Validate(), test.Err, "Validate(%s)", test.Name)

		name, err := test.Packet.CNAME()
		assert.ErrorIsf(t, err, test.Err, "CNAME(%s)", test.Name)
		assert.Equalf(t, test.Text, name, "CNAME(%s)", test.Name)
	}
}

func TestCompoundPacketRoundTrip(t *testing.T) {
	cname := NewCNAMESourceDescription(1234, "cname")

	for _, test := range []struct {
		Name   string
		Packet CompoundPacket
		Err    error
	}{
		{
			Name: "bye",
			Packet: CompoundPacket{
				&ReceiverReport{},
				cname,
				&Goodbye{
					Sources: []uint32{1234},
				},
			},
		},
		{
			Name: "no cname",
			Packet: CompoundPacket{
				&ReceiverReport{},
			},
			Err: errMissingCNAME,
		},
	} {
		data, err := test.Packet.Marshal()
		assert.ErrorIsf(t, err, test.Err, "Marshal(%v)", test.Name)
		if err != nil {
			continue
		}

		var c CompoundPacket
		assert.NoErrorf(t, c.Unmarshal(data), "Unmarshal(%v)", test.Name)

		data2, err := c.Marshal()
		assert.NoErrorf(t, err, "Marshal(%v)", test.Name)
		assert.Equalf(t, data, data2, "Marshal(%v) mismatch", test.Name)
	}
}
