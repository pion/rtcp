// SPDX-FileCopyrightText: 2025 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Packet = (*TMMBN)(nil) // assert is a Packet

func TestTMMBNMarshal(t *testing.T) {
	assert := assert.New(t)

	input := TMMBN{
		SenderSSRC: 1,
		Entries: []TMMBNEntry{
			{
				MediaSSRC: 1215622422,
				Bitrate:   8927168.0,
			},
		},
	}

	// Expected packet structure:
	// Header: V=2, P=0, FMT=4, PT=205, Length=4
	// SenderSSRC: 0x00000001
	// MediaSSRC: 0x00000000 (always 0 per RFC 5104)
	// FCI Entry:
	//   - SSRC: 0x48746ED6 (1215622422)
	//   - Bitrate: exp=6, mantissa=139487 (0x0220DF) -> 0x1A20DF
	//   - Overhead: 0x00
	expected := []byte{132, 205, 0, 4, 0, 0, 0, 1, 0, 0, 0, 0, 72, 116, 237, 22, 26, 32, 223, 0}

	output, err := input.Marshal()
	assert.NoError(err)
	assert.Equal(expected, output)
}

func TestTMMBNUnmarshal(t *testing.T) {
	assert := assert.New(t)

	// Real TMMBN packet with bitrate 8927168 (8.9 Mb/s)
	input := []byte{132, 205, 0, 4, 0, 0, 0, 1, 0, 0, 0, 0, 72, 116, 237, 22, 26, 32, 223, 0}

	// mantissa = []byte{26 & 3, 32, 223} = []byte{2, 32, 223} = 139487
	// exp = 26 >> 2 = 6
	// bitrate = 139487 * 2^6 = 139487 * 64 = 8927168 = 8.9 Mb/s
	expected := TMMBN{
		SenderSSRC: 1,
		Entries: []TMMBNEntry{
			{
				MediaSSRC: 1215622422,
				Bitrate:   8927168,
			},
		},
	}

	packet := TMMBN{}
	err := packet.Unmarshal(input)
	assert.NoError(err)
	assert.Equal(expected, packet)
}

func TestTMMBNTruncate(t *testing.T) {
	assert := assert.New(t)

	input := []byte{132, 205, 0, 4, 0, 0, 0, 1, 0, 0, 0, 0, 72, 116, 237, 22, 26, 32, 223, 0}

	// Make sure that we're interpreting the bitrate correctly.
	// For the above example, we have:

	// mantissa = 139487
	// exp = 6
	// bitrate = 8927168

	packet := TMMBN{}
	err := packet.Unmarshal(input)
	assert.NoError(err)
	assert.Equal(float32(8927168), packet.Entries[0].Bitrate)

	// Just verify marshal produces the same input.
	output, err := packet.Marshal()
	assert.NoError(err)
	assert.Equal(input, output)

	// If we subtract the bitrate by 1, we'll round down a lower mantissa
	packet.Entries[0].Bitrate--

	// bitrate = 8927167
	// mantissa = 139486
	// exp = 6

	output, err = packet.Marshal()
	assert.NoError(err)
	assert.NotEqual(input, output)

	// Which if we actually unmarshal again, we'll find that it's actually decreased by 64 (which is 2^exp)
	// mantissa = 139486
	// exp = 6
	// bitrate = 8927104

	err = packet.Unmarshal(output)
	assert.NoError(err)
	assert.Equal(float32(8927104), packet.Entries[0].Bitrate)
}

func TestTMMBNOverflow(t *testing.T) {
	assert := assert.New(t)

	// Marshal a packet with the maximum possible bitrate.
	packet := TMMBN{
		Entries: []TMMBNEntry{
			{
				Bitrate: math.MaxFloat32,
			},
		},
	}

	// mantissa = 262143 = 0x3FFFF
	// exp = 63

	expected := []byte{132, 205, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 0}

	output, err := packet.Marshal()
	assert.NoError(err)
	assert.Equal(expected, output)

	// mantissa = 262143
	// exp = 63
	// bitrate = 0xFFFFC00000000000

	err = packet.Unmarshal(output)
	assert.NoError(err)
	assert.Equal(math.Float32frombits(0x67FFFFC0), packet.Entries[0].Bitrate)

	// Make sure we marshal to the same result again.
	output, err = packet.Marshal()
	assert.NoError(err)
	assert.Equal(expected, output)

	// Finally, try unmarshalling one number higher than we used to be able to handle.
	input := []byte{132, 205, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 188, 0, 0, 0}
	err = packet.Unmarshal(input)
	assert.NoError(err)
	assert.Equal(math.Float32frombits(0x62800000), packet.Entries[0].Bitrate)
}

func TestTMMBNMultipleEntries(t *testing.T) {
	assert := assert.New(t)

	input := TMMBN{
		SenderSSRC: 12345,
		Entries: []TMMBNEntry{
			{
				MediaSSRC: 1000,
				Bitrate:   1000000.0,
			},
			{
				MediaSSRC: 2000,
				Bitrate:   2000000.0,
			},
		},
	}

	output, err := input.Marshal()
	assert.NoError(err)

	packet := TMMBN{}
	err = packet.Unmarshal(output)
	assert.NoError(err)

	assert.Equal(input.SenderSSRC, packet.SenderSSRC)
	assert.Equal(len(input.Entries), len(packet.Entries))

	for i := range input.Entries {
		assert.Equal(input.Entries[i].MediaSSRC, packet.Entries[i].MediaSSRC)
		// Allow small floating point differences due to encoding/decoding
		assert.InDelta(input.Entries[i].Bitrate, packet.Entries[i].Bitrate, 1000)
	}
}

func TestTMMBNDestinationSSRC(t *testing.T) {
	assert := assert.New(t)

	packet := TMMBN{
		Entries: []TMMBNEntry{
			{MediaSSRC: 1000},
			{MediaSSRC: 2000},
			{MediaSSRC: 3000},
		},
	}

	ssrcs := packet.DestinationSSRC()
	assert.Equal([]uint32{1000, 2000, 3000}, ssrcs)
}

func TestTMMBNMarshalSize(t *testing.T) {
	assert := assert.New(t)

	// Test with no entries
	packet := TMMBN{}
	assert.Equal(12, packet.MarshalSize())

	// Test with one entry
	packet.Entries = []TMMBNEntry{{}}
	assert.Equal(20, packet.MarshalSize())

	// Test with multiple entries
	packet.Entries = []TMMBNEntry{{}, {}, {}}
	assert.Equal(36, packet.MarshalSize())
}

func TestTMMBNHeader(t *testing.T) {
	assert := assert.New(t)

	packet := TMMBN{
		SenderSSRC: 1,
		Entries: []TMMBNEntry{
			{MediaSSRC: 1000, Bitrate: 1000000},
		},
	}

	header := packet.Header()
	assert.Equal(FormatTMMBN, int(header.Count))
	assert.Equal(TypeTransportSpecificFeedback, header.Type)
	assert.Equal(4, int(header.Length))
}

func TestTMMBNString(t *testing.T) {
	assert := assert.New(t)

	packet := TMMBN{
		SenderSSRC: 0x12345678,
		Entries: []TMMBNEntry{
			{
				MediaSSRC: 0xABCDEF00,
				Bitrate:   8927168.0,
			},
		},
	}

	str := packet.String()
	assert.Contains(str, "TMMBN")
	assert.Contains(str, "12345678")
	assert.Contains(str, "abcdef00")
}

func TestTMMBNUnmarshalErrors(t *testing.T) {
	assert := assert.New(t)

	// Test packet too short
	packet := TMMBN{}
	err := packet.Unmarshal([]byte{1, 2, 3})
	assert.Error(err)

	// Test wrong packet type
	wrongType := []byte{132, 200, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	err = packet.Unmarshal(wrongType)
	assert.Error(err)

	// Test wrong format
	wrongFormat := []byte{132, 205, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	wrongFormat[0] = 135 // Change FMT to 7
	err = packet.Unmarshal(wrongFormat)
	assert.Error(err)
}
