// SPDX-FileCopyrightText: 2025 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// TMMBN represents a Temporary Maximum Media Stream Bit Rate Notification packet
// as defined in RFC 5104, section 4.2.2.
type TMMBN struct {
	// SSRC of the sender
	SenderSSRC uint32

	// List of TMMBN entries
	Entries []TMMBNEntry
}

// TMMBNEntry represents a single entry in TMMBN packet
type TMMBNEntry struct {
	// SSRC of media source this entry applies to
	MediaSSRC uint32

	// Estimated maximum bitrate
	Bitrate float32
}

// Marshal encodes the TMMBN packet in binary format
func (p TMMBN) Marshal() ([]byte, error) {
	/*
		TMMBN packet format (RFC 5104):
		 0                   1                   2                   3
		 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|V=2|P|  FMT=4  |   PT = 205    |          length               |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                  SSRC of RTCP packet sender                   |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                  SSRC of media source                         |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                           FCI SSRC                            |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		| MxTBR Exp |  MxTBR Mantissa                 |Measured Overhead|
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|  ...                                                          |
	*/

	packetSize := p.MarshalSize()
	rawPacket := make([]byte, packetSize)

	header := p.Header()
	headerBuf, err := header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(rawPacket, headerBuf)

	body := rawPacket[headerLength:]
	binary.BigEndian.PutUint32(body, p.SenderSSRC)
	// Media SSRC is always 0
	// https://www.rfc-editor.org/rfc/rfc5104.html#section-4.2.1.2
	binary.BigEndian.PutUint32(body[ssrcLength:], 0)

	// Write each FCI entry
	for i, entry := range p.Entries {
		offset := ssrcLength*2 + i*(2*ssrcLength)
		binary.BigEndian.PutUint32(body[offset:], entry.MediaSSRC)

		err = putBitrate(entry.Bitrate, body[offset+ssrcLength:])
		if err != nil {
			return nil, err
		}
	}

	return rawPacket, nil
}

// Unmarshal decodes the TMMBN packet from binary data
func (p *TMMBN) Unmarshal(rawPacket []byte) error {
	if len(rawPacket) < headerLength+ssrcLength*2 {
		return errPacketTooShort
	}

	var header Header
	if err := header.Unmarshal(rawPacket); err != nil {
		return err
	}

	expectedSize := int((header.Length + 1) * 4)
	if len(rawPacket) < expectedSize {
		return errBadLength
	}

	if header.Type != TypeTransportSpecificFeedback || header.Count != FormatTMMBN {
		return errWrongType
	}

	body := rawPacket[headerLength:]
	p.SenderSSRC = binary.BigEndian.Uint32(body)

	entryCount := int((header.Length - 2) / 2)
	p.Entries = make([]TMMBNEntry, entryCount)

	for i := 0; i < entryCount; i++ {
		offset := ssrcLength*2 + i*(2*ssrcLength)
		entry := &p.Entries[i]
		entry.MediaSSRC = binary.BigEndian.Uint32(body[offset:])
		entry.Bitrate = loadBitrate(body[offset+ssrcLength:])
	}

	return nil
}

// MarshalSize returns the size of the packet when marshaled
func (p *TMMBN) MarshalSize() int {
	return headerLength + ssrcLength*2 + len(p.Entries)*(2*ssrcLength)
}

func (p *TMMBN) Header() Header {
	return Header{
		Count:  FormatTMMBN,
		Type:   TypeTransportSpecificFeedback,
		Length: uint16((p.MarshalSize() / 4) - 1), //nolint:gosec // G115
	}
}

// String prints the TMMBN packet in a human-readable format.
func (p *TMMBN) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TMMBN from %x:\n", p.SenderSSRC))
	for i, entry := range p.Entries {
		unit := bitrateUnit(entry.Bitrate)
		sb.WriteString(fmt.Sprintf(" entry %d: media=%x, bitrate=%.2f %s/s\n", i, entry.MediaSSRC, entry.Bitrate, unit))
	}
	return sb.String()
}

// DestinationSSRC returns SSRCs this packet applies to
func (p *TMMBN) DestinationSSRC() []uint32 {
	ssrcs := make([]uint32, len(p.Entries))
	for i, entry := range p.Entries {
		ssrcs[i] = entry.MediaSSRC
	}
	return ssrcs
}
