// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import "math"

// getPadding Returns the padding required to make the length a multiple of 4.
func getPadding(packetLen int) int {
	if packetLen%4 == 0 {
		return 0
	}

	return 4 - (packetLen % 4)
}

// setNBitsOfUint16 will truncate the value to size, left-shift to startIndex position and set.
func setNBitsOfUint16(src, size, startIndex, val uint16) (uint16, error) {
	if startIndex+size > 16 {
		return 0, errInvalidSizeOrStartIndex
	}

	// truncate val to size bits
	val &= (1 << size) - 1

	return src | (val << (16 - size - startIndex)), nil
}

// appendBit32 will left-shift and append n bits of val.
func appendNBitsToUint32(src, n, val uint32) uint32 {
	return (src << n) | (val & (0xFFFFFFFF >> (32 - n)))
}

// getNBit get n bits from 1 byte, begin with a position.
func getNBitsFromByte(b byte, begin, n uint16) uint16 {
	endShift := 8 - (begin + n)
	mask := (0xFF >> begin) & uint8(0xFF<<endShift)

	return uint16(b&mask) >> endShift
}

// get24BitFromBytes get 24bits from `[3]byte` slice.
func get24BitsFromBytes(b []byte) uint32 {
	return uint32(b[0])<<16 + uint32(b[1])<<8 + uint32(b[2])
}

func putBitrate(bitrate float32, buf []byte) (err error) {
	const bitratemax = 0x3FFFFp+63
	if bitrate >= bitratemax {
		bitrate = bitratemax
	}

	if bitrate < 0 {
		return errInvalidBitrate
	}

	exp := 0

	for bitrate >= (1 << 18) {
		bitrate /= 2.0
		exp++
	}

	if exp >= (1 << 6) {
		return errInvalidBitrate
	}

	mantissa := uint(math.Floor(float64(bitrate)))

	// We can't quite use the binary package because
	// a) it's a uint24 and b) the exponent is only 6-bits
	// Just trust me; this is big-endian encoding.
	buf[0] = byte(exp<<2) | byte(mantissa>>16)
	buf[1] = byte(mantissa >> 8)
	buf[2] = byte(mantissa)
	return nil
}

func loadBitrate(buf []byte) float32 {
	const mantissamax = 0x7FFFFF
	// Get the 6-bit exponent value.
	exp := buf[0] >> 2
	exp += 127 // bias for IEEE754
	exp += 23  // IEEE754 biases the decimal to the left, abs-send-time biases it to the right

	// The remaining 2-bits plus the next 16-bits are the mantissa.
	mantissa := uint32(buf[0]&3)<<16 | uint32(buf[1])<<8 | uint32(buf[2])

	if mantissa != 0 {
		// ieee754 requires an implicit leading bit
		for (mantissa & (mantissamax + 1)) == 0 {
			exp--
			mantissa *= 2
		}
	}

	// bitrate = mantissa * 2^exp
	bitrate := math.Float32frombits((uint32(exp) << 23) | (mantissa & mantissamax))
	return bitrate
}

func bitrateUnit(bitrate float32) string {
	// Keep a table of powers to units for fast conversion.
	bitUnits := []string{"b", "Kb", "Mb", "Gb", "Tb", "Pb", "Eb"}

	// Do some unit conversions because b/s is far too difficult to read.
	powers := 0

	// Keep dividing the bitrate until it's under 1000
	for bitrate >= 1000.0 && powers < len(bitUnits) {
		bitrate /= 1000.0
		powers++
	}

	return bitUnits[powers] //nolint:gosec // powers is bounded by loop condition
}
