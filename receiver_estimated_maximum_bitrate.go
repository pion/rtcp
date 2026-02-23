// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ReceiverEstimatedMaximumBitrate contains the receiver's estimated maximum bitrate.
// see: https://tools.ietf.org/html/draft-alvestrand-rmcat-remb-03
type ReceiverEstimatedMaximumBitrate struct {
	// SSRC of sender
	SenderSSRC uint32

	// Estimated maximum bitrate
	Bitrate float32

	// SSRC entries which this packet applies to
	SSRCs []uint32
}

// Marshal serializes the packet and returns a byte slice.
func (p ReceiverEstimatedMaximumBitrate) Marshal() (buf []byte, err error) {
	// Allocate a buffer of the exact output size.
	buf = make([]byte, p.MarshalSize())

	// Write to our buffer.
	n, err := p.MarshalTo(buf)
	if err != nil {
		return nil, err
	}

	// This will always be true but just to be safe.
	if n != len(buf) {
		return nil, errWrongMarshalSize
	}

	return buf, nil
}

// MarshalSize returns the size of the packet once marshaled.
func (p ReceiverEstimatedMaximumBitrate) MarshalSize() int {
	return 20 + 4*len(p.SSRCs)
}

// MarshalTo serializes the packet to the given byte slice.
func (p ReceiverEstimatedMaximumBitrate) MarshalTo(buf []byte) (n int, err error) {
	/*
	    0                   1                   2                   3
	    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |V=2|P| FMT=15  |   PT=206      |             length            |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |                  SSRC of packet sender                        |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |                  SSRC of media source                         |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |  Unique identifier 'R' 'E' 'M' 'B'                            |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |  Num SSRC     | BR Exp    |  BR Mantissa                      |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |   SSRC feedback                                               |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |  ...                                                          |
	*/

	size := p.MarshalSize()
	if len(buf) < size {
		return 0, errPacketTooShort
	}

	buf[0] = 143 // v=2, p=0, fmt=15
	buf[1] = 206

	// Length of this packet in 32-bit words minus one.
	length := uint16((p.MarshalSize() / 4) - 1) //nolint:gosec // G115
	binary.BigEndian.PutUint16(buf[2:4], length)

	binary.BigEndian.PutUint32(buf[4:8], p.SenderSSRC)
	binary.BigEndian.PutUint32(buf[8:12], 0) // always zero

	// ALL HAIL REMB
	buf[12] = 'R'
	buf[13] = 'E'
	buf[14] = 'M'
	buf[15] = 'B'

	// Write the length of the ssrcs to follow at the end
	buf[16] = byte(len(p.SSRCs))

	err = putBitrate(p.Bitrate, buf[17:20])
	if err != nil {
		return 0, err
	}

	// Write the SSRCs at the very end.
	n = 20
	for _, ssrc := range p.SSRCs {
		binary.BigEndian.PutUint32(buf[n:n+4], ssrc)
		n += 4
	}

	return n, nil
}

// Unmarshal reads a REMB packet from the given byte slice.
//
//nolint:cyclop
func (p *ReceiverEstimatedMaximumBitrate) Unmarshal(buf []byte) (err error) {
	/*
	    0                   1                   2                   3
	    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |V=2|P| FMT=15  |   PT=206      |             length            |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |                  SSRC of packet sender                        |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |                  SSRC of media source                         |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |  Unique identifier 'R' 'E' 'M' 'B'                            |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |  Num SSRC     | BR Exp    |  BR Mantissa                      |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |   SSRC feedback                                               |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |  ...                                                          |
	*/

	// 20 bytes is the size of the packet with no SSRCs
	if len(buf) < 20 {
		return errPacketTooShort
	}

	// version  must be 2
	version := buf[0] >> 6
	if version != 2 {
		return fmt.Errorf("%w expected(2) actual(%d)", errBadVersion, version)
	}

	// padding must be unset
	padding := (buf[0] >> 5) & 1
	if padding != 0 {
		return fmt.Errorf("%w expected(0) actual(%d)", errWrongPadding, padding)
	}

	// fmt must be 15
	fmtVal := buf[0] & 31
	if fmtVal != 15 {
		return fmt.Errorf("%w expected(15) actual(%d)", errWrongFeedbackType, fmtVal)
	}

	// Must be payload specific feedback
	if buf[1] != 206 {
		return fmt.Errorf("%w expected(206) actual(%d)", errWrongPayloadType, buf[1])
	}

	// length is the number of 32-bit words, minus 1
	length := binary.BigEndian.Uint16(buf[2:4])
	size := int((length + 1) * 4)

	// There's not way this could be legit
	if size < 20 {
		return errHeaderTooSmall
	}

	// Make sure the buffer is large enough.
	if len(buf) < size {
		return errPacketTooShort
	}

	// The sender SSRC is 32-bits
	p.SenderSSRC = binary.BigEndian.Uint32(buf[4:8])

	// The destination SSRC must be 0
	media := binary.BigEndian.Uint32(buf[8:12])
	if media != 0 {
		return errSSRCMustBeZero
	}

	// REMB rules all around me
	if !bytes.Equal(buf[12:16], []byte{'R', 'E', 'M', 'B'}) {
		return errMissingREMBidentifier
	}

	// The next byte is the number of SSRC entries at the end.
	num := int(buf[16])

	// Now we know the expected size, make sure they match.
	if size != 20+4*num {
		return errSSRCNumAndLengthMismatch
	}

	p.Bitrate = loadBitrate(buf[17:20])

	// Clear any existing SSRCs
	p.SSRCs = nil

	// Loop over and parse the SSRC entires at the end.
	// We already verified that size == num * 4
	for n := 20; n < size; n += 4 {
		ssrc := binary.BigEndian.Uint32(buf[n : n+4])
		p.SSRCs = append(p.SSRCs, ssrc)
	}

	return nil
}

// Header returns the Header associated with this packet.
func (p *ReceiverEstimatedMaximumBitrate) Header() Header {
	return Header{
		Count:  FormatREMB,
		Type:   TypePayloadSpecificFeedback,
		Length: uint16((p.MarshalSize() / 4) - 1), //nolint:gosec // G115
	}
}

// String prints the REMB packet in a human-readable format.
func (p *ReceiverEstimatedMaximumBitrate) String() string {
	unit := bitrateUnit(p.Bitrate)
	return fmt.Sprintf("ReceiverEstimatedMaximumBitrate %x %.2f %s/s", p.SenderSSRC, p.Bitrate, unit)
}

// DestinationSSRC returns an array of SSRC values that this packet refers to.
func (p *ReceiverEstimatedMaximumBitrate) DestinationSSRC() []uint32 {
	return p.SSRCs
}
