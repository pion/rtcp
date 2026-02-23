// SPDX-FileCopyrightText: 2025 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// TMMBR represents a Temporary Maximum Media Stream Bit Rate Request packet
// as defined in RFC 5104, section 4.2.1.
type TMMBR struct {
	// SSRC of the sender
	SenderSSRC uint32

	// List of TMMBR entries
	Entries []TMMBREntry
}

// TMMBREntry represents a single entry in TMMBR packet
type TMMBREntry struct {
	// SSRC of media source this entry applies to
	MediaSSRC uint32

	// Estimated maximum bitrate
	Bitrate float32
}

// Marshal encodes the TMMBR packet in binary format
func (p TMMBR) Marshal() ([]byte, error) {
	/*
		TMMBR packet format (RFC 5104):
		 0                   1                   2                   3
		 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|V=2|P|  FMT=3  |   PT = 205    |          length               |
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

// Unmarshal decodes the TMMBR packet from binary data
func (p *TMMBR) Unmarshal(rawPacket []byte) error {
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

	if header.Type != TypeTransportSpecificFeedback || header.Count != FormatTMMBR {
		return errWrongType
	}

	body := rawPacket[headerLength:]
	p.SenderSSRC = binary.BigEndian.Uint32(body)

	entryCount := int((header.Length - 2) / 2)
	p.Entries = make([]TMMBREntry, entryCount)

	for i := 0; i < entryCount; i++ {
		offset := ssrcLength*2 + i*(2*ssrcLength)
		entry := &p.Entries[i]
		entry.MediaSSRC = binary.BigEndian.Uint32(body[offset:])
		entry.Bitrate = loadBitrate(body[offset+ssrcLength:])
	}

	return nil
}

// MarshalSize returns the size of the packet when marshaled
func (p *TMMBR) MarshalSize() int {
	return headerLength + ssrcLength*2 + len(p.Entries)*(2*ssrcLength)
}

func (p *TMMBR) Header() Header {
	return Header{
		Count:  FormatTMMBR,
		Type:   TypeTransportSpecificFeedback,
		Length: uint16((p.MarshalSize() / 4) - 1), //nolint:gosec // G115
	}
}

// String prints the TMMBR packet in a human-readable format.
func (p *TMMBR) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TMMBR from %x:\n", p.SenderSSRC))
	for i, entry := range p.Entries {
		unit := bitrateUnit(entry.Bitrate)
		sb.WriteString(fmt.Sprintf(" entry %d: media=%x, bitrate=%.2f %s/s\n", i, entry.MediaSSRC, entry.Bitrate, unit))
	}
	return sb.String()
}

// DestinationSSRC returns SSRCs this packet applies to
func (p *TMMBR) DestinationSSRC() []uint32 {
	ssrcs := make([]uint32, len(p.Entries))
	for i, entry := range p.Entries {
		ssrcs[i] = entry.MediaSSRC
	}
	return ssrcs
}
