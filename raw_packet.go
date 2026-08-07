// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"encoding/binary"
	"fmt"
)

// RawPacket represents an unparsed RTCP packet. It's returned by Unmarshal when
// a packet with an unknown type is encountered.
type RawPacket []byte

// Marshal encodes the packet in binary.
func (r RawPacket) Marshal() ([]byte, error) {
	return r, nil
}

// Unmarshal decodes the packet from binary.
func (r *RawPacket) Unmarshal(b []byte) error {
	if len(b) < (headerLength) {
		return errPacketTooShort
	}
	*r = b

	var h Header

	return h.Unmarshal(b)
}

// Header returns the Header associated with this packet.
func (r RawPacket) Header() Header {
	var h Header
	if err := h.Unmarshal(r); err != nil {
		return Header{}
	}

	return h
}

// DestinationSSRC returns an array of SSRC values that this packet refers to.
func (r *RawPacket) DestinationSSRC() []uint32 {
	return []uint32{}
}

func (r RawPacket) String() string {
	out := fmt.Sprintf("RawPacket: %v", ([]byte)(r))

	return out
}

// MarshalSize returns the size of the packet once marshaled.
func (r RawPacket) MarshalSize() int {
	return len(r)
}

// UnmarshalRaw takes an entire udp datagram (which may consist of multiple
// RTCP packets) and returns the raw packets it contains, without unmarshaling
// their contents.
//
// The returned packets are subslices of data, not copies.
func UnmarshalRaw(data []byte) ([]RawPacket, error) {
	return AppendRawPackets(nil, data)
}

// AppendRawPackets appends the raw packets of the datagram data, as produced
// by UnmarshalRaw, to dst and returns the extended slice. On error dst is
// returned unchanged.
func AppendRawPackets(dst []RawPacket, data []byte) ([]RawPacket, error) {
	orig := dst
	found := 0
	for rest := data; len(rest) > 0; {
		var header Header
		if err := header.Unmarshal(rest); err != nil {
			return orig, err
		}
		size := (int(header.Length) + 1) * 4
		if size > len(rest) {
			return orig, errPacketTooShort
		}
		dst = append(dst, RawPacket(rest[:size:size]))
		rest = rest[size:]
		found++
	}
	if found == 0 {
		return orig, errInvalidHeader
	}

	return dst, nil
}

// ParseDestinationSSRC parses the destination SSRCs of the packet out of its
// raw bytes, appends them to dst, and returns the extended slice. It reports
// the SSRCs that Unmarshal followed by Packet.DestinationSSRC would for every
// packet type known to this package; unknown packet types have none. r must
// be a single RTCP packet, as produced by UnmarshalRaw.
//
// Only the structure needed to locate the SSRCs is validated; a packet
// accepted by this method may still fail a full Unmarshal.
//
// ParseDestinationSSRC will only allocate in order to resize dst as needed;
// the parsing itself is allocation-free.
//
//nolint:cyclop
func (r RawPacket) ParseDestinationSSRC(dst []uint32) ([]uint32, error) {
	pkt := []byte(r)
	var header Header
	if err := header.Unmarshal(pkt); err != nil {
		return dst, err
	}
	if size := (int(header.Length) + 1) * 4; size != len(pkt) {
		return dst, errBadLength
	}

	switch header.Type {
	case TypeSenderReport:
		return appendSenderReportSSRCs(dst, header, pkt)
	case TypeReceiverReport:
		return appendReceiverReportSSRCs(dst, header, pkt)
	case TypeSourceDescription:
		return appendSourceDescriptionSSRCs(dst, header, pkt)
	case TypeGoodbye:
		return appendGoodbyeSSRCs(dst, header, pkt)
	case TypeApplicationDefined:
		return appendApplicationDefinedSSRCs(dst, header, pkt)
	case TypeTransportSpecificFeedback:
		switch header.Count {
		case FormatRRR:
			var rrr RapidResynchronizationRequest
			if err := rrr.Unmarshal(pkt); err != nil {
				return dst, err
			}

			return append(dst, rrr.MediaSSRC), nil
		case FormatTLN, FormatTCC:
			return appendMediaSSRC(dst, header, pkt)
		case FormatCCFB:
			return appendCCFBSSRCs(dst, header, pkt)
		}
	case TypePayloadSpecificFeedback:
		switch header.Count {
		case FormatPLI:
			var pli PictureLossIndication
			if err := pli.Unmarshal(pkt); err != nil {
				return dst, err
			}

			return append(dst, pli.MediaSSRC), nil
		case FormatSLI:
			return appendMediaSSRC(dst, header, pkt)
		case FormatFIR:
			return appendFIRSSRCs(dst, header, pkt)
		case FormatREMB:
			return appendREMBSSRCs(dst, header, pkt)
		}
	case TypeExtendedReport:
		return appendXRSSRCs(dst, header, pkt)
	}

	return dst, nil
}

// appendMediaSSRC handles the feedback packets whose only destination is the
// media SSRC in the fixed part of the packet and whose full Unmarshal
// allocates (SLI, NACK, TWCC). They share the common packet format for
// feedback messages:
//
// https://tools.ietf.org/html/rfc4585#section-6.1
//
//	 0                   1                   2                   3
//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|V=2|P|   FMT   |       PT      |          length               |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                  SSRC of packet sender                        |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                  SSRC of media source                         |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	:            Feedback Control Information (FCI)                 :
//	:                                                               :
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
func appendMediaSSRC(dst []uint32, _ Header, pkt []byte) ([]uint32, error) {
	if len(pkt) < (headerLength + (ssrcLength * 2)) {
		return dst, errPacketTooShort
	}

	return append(dst, binary.BigEndian.Uint32(pkt[headerLength+ssrcLength:])), nil
}
