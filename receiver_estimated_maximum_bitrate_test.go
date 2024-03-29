// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Packet = (*ReceiverEstimatedMaximumBitrate)(nil) // assert is a Packet

func TestReceiverEstimatedMaximumBitrateMarshal(t *testing.T) {
	assert := assert.New(t)

	input := ReceiverEstimatedMaximumBitrate{
		SenderSSRC: 1,
		Bitrate:    8927168.0,
		SSRCs:      []uint32{1215622422},
	}

	expected := []byte{143, 206, 0, 5, 0, 0, 0, 1, 0, 0, 0, 0, 82, 69, 77, 66, 1, 26, 32, 223, 72, 116, 237, 22}

	output, err := input.Marshal()
	assert.NoError(err)
	assert.Equal(expected, output)
}

func TestReceiverEstimatedMaximumBitrateUnmarshal(t *testing.T) {
	assert := assert.New(t)

	// Real data sent by Chrome while watching a 6Mb/s stream
	input := []byte{143, 206, 0, 5, 0, 0, 0, 1, 0, 0, 0, 0, 82, 69, 77, 66, 1, 26, 32, 223, 72, 116, 237, 22}

	// mantissa = []byte{26 & 3, 32, 223} = []byte{2, 32, 223} = 139487
	// exp = 26 >> 2 = 6
	// bitrate = 139487 * 2^6 = 139487 * 64 = 8927168 = 8.9 Mb/s
	expected := ReceiverEstimatedMaximumBitrate{
		SenderSSRC: 1,
		Bitrate:    8927168,
		SSRCs:      []uint32{1215622422},
	}

	packet := ReceiverEstimatedMaximumBitrate{}
	err := packet.Unmarshal(input)
	assert.NoError(err)
	assert.Equal(expected, packet)
}

func TestReceiverEstimatedMaximumBitrateTruncate(t *testing.T) {
	assert := assert.New(t)

	input := []byte{143, 206, 0, 5, 0, 0, 0, 1, 0, 0, 0, 0, 82, 69, 77, 66, 1, 26, 32, 223, 72, 116, 237, 22}

	// Make sure that we're interpreting the bitrate correctly.
	// For the above example, we have:

	// mantissa = 139487
	// exp = 6
	// bitrate = 8927168

	packet := ReceiverEstimatedMaximumBitrate{}
	err := packet.Unmarshal(input)
	assert.NoError(err)
	assert.Equal(float32(8927168), packet.Bitrate)

	// Just verify marshal produces the same input.
	output, err := packet.Marshal()
	assert.NoError(err)
	assert.Equal(input, output)

	// If we subtract the bitrate by 1, we'll round down a lower mantissa
	packet.Bitrate--

	// bitrate = 8927167
	// mantissa = 139486
	// exp = 6

	output, err = packet.Marshal()
	assert.NoError(err)
	assert.NotEqual(input, output)

	// Which if we actually unmarshal again, we'll find that it's actually decreased by 63 (which is exp)
	// mantissa = 139486
	// exp = 6
	// bitrate = 8927104

	err = packet.Unmarshal(output)
	assert.NoError(err)
	assert.Equal(float32(8927104), packet.Bitrate)
}

func TestReceiverEstimatedMaximumBitrateOverflow(t *testing.T) {
	assert := assert.New(t)

	// Marshal a packet with the maximum possible bitrate.
	packet := ReceiverEstimatedMaximumBitrate{
		Bitrate: math.MaxFloat32,
	}

	// mantissa = 262143 = 0x3FFFF
	// exp = 63

	expected := []byte{143, 206, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 82, 69, 77, 66, 0, 255, 255, 255}

	output, err := packet.Marshal()
	assert.NoError(err)
	assert.Equal(expected, output)

	// mantissa = 262143
	// exp = 63
	// bitrate = 0xFFFFC00000000000

	err = packet.Unmarshal(output)
	assert.NoError(err)
	assert.Equal(math.Float32frombits(0x67FFFFC0), packet.Bitrate)

	// Make sure we marshal to the same result again.
	output, err = packet.Marshal()
	assert.NoError(err)
	assert.Equal(expected, output)

	// Finally, try unmarshalling one number higher than we used to be able to handle.
	input := []byte{143, 206, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 82, 69, 77, 66, 0, 188, 0, 0}
	err = packet.Unmarshal(input)
	assert.NoError(err)
	assert.Equal(math.Float32frombits(0x62800000), packet.Bitrate)
}
