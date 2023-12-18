// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import "encoding/binary"

// A ReceptionReport block conveys statistics on the reception of RTP packets
// from a single synchronization source.
type ReceptionReport struct {
	// The SSRC identifier of the source to which the information in this
	// reception report block pertains.
	SSRC uint32
	// The fraction of RTP data packets from source SSRC lost since the
	// previous SR or RR packet was sent, expressed as a fixed point
	// number with the binary point at the left edge of the field.
	FractionLost uint8
	// The total number of RTP data packets from source SSRC that have
	// been lost since the beginning of reception.
	// https://www.ietf.org/rfc/rfc3550.txt Section 6.4 and appendix A.3
	// The total number of packets lost is a signed number and can be
	// negative
	TotalLost int32
	// The low 16 bits contain the highest sequence number received in an
	// RTP data packet from source SSRC, and the most significant 16
	// bits extend that sequence number with the corresponding count of
	// sequence number cycles.
	LastSequenceNumber uint32
	// An estimate of the statistical variance of the RTP data packet
	// interarrival time, measured in timestamp units and expressed as an
	// unsigned integer.
	Jitter uint32
	// The middle 32 bits out of 64 in the NTP timestamp received as part of
	// the most recent RTCP sender report (SR) packet from source SSRC. If no
	// SR has been received yet, the field is set to zero.
	LastSenderReport uint32
	// The delay, expressed in units of 1/65536 seconds, between receiving the
	// last SR packet from source SSRC and sending this reception report block.
	// If no SR packet has been received yet from SSRC, the field is set to zero.
	Delay uint32
}

const (
	receptionReportLength = 24
	fractionLostOffset    = 4
	totalLostOffset       = 5
	lastSeqOffset         = 8
	jitterOffset          = 12
	lastSROffset          = 16
	delayOffset           = 20
)

// Marshal encodes the ReceptionReport in binary
func (r ReceptionReport) Marshal() ([]byte, error) {
	/*
	 *  0                   1                   2                   3
	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	 * +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
	 * |                              SSRC                             |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * | fraction lost |       cumulative number of packets lost       |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |           extended highest sequence number received           |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                      interarrival jitter                      |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                         last SR (LSR)                         |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                   delay since last SR (DLSR)                  |
	 * +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
	 */

	rawPacket := make([]byte, receptionReportLength)

	binary.BigEndian.PutUint32(rawPacket, r.SSRC)

	rawPacket[fractionLostOffset] = r.FractionLost

	// pack TotalLost into 24 bits
	// we first convert signed integer to unsigned before using bit operators
	uTotalLost := uint32(r.TotalLost)
	if uTotalLost&0xff800000 != 0xff800000 && uTotalLost&0xff800000 != 0 {
		return nil, errInvalidTotalLost
	}
	// Convert int32 to int24
	uTotalLost = uTotalLost&0x80000000>>8 | uTotalLost&0x007fffff
	tlBytes := rawPacket[totalLostOffset:]
	tlBytes[0] = byte(uTotalLost >> 16)
	tlBytes[1] = byte(uTotalLost >> 8)
	tlBytes[2] = byte(uTotalLost)

	binary.BigEndian.PutUint32(rawPacket[lastSeqOffset:], r.LastSequenceNumber)
	binary.BigEndian.PutUint32(rawPacket[jitterOffset:], r.Jitter)
	binary.BigEndian.PutUint32(rawPacket[lastSROffset:], r.LastSenderReport)
	binary.BigEndian.PutUint32(rawPacket[delayOffset:], r.Delay)

	return rawPacket, nil
}

// Unmarshal decodes the ReceptionReport from binary
func (r *ReceptionReport) Unmarshal(rawPacket []byte) error {
	if len(rawPacket) < receptionReportLength {
		return errPacketTooShort
	}

	/*
	 *  0                   1                   2                   3
	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	 * +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
	 * |                              SSRC                             |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * | fraction lost |       cumulative number of packets lost       |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |           extended highest sequence number received           |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                      interarrival jitter                      |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                         last SR (LSR)                         |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                   delay since last SR (DLSR)                  |
	 * +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
	 */

	r.SSRC = binary.BigEndian.Uint32(rawPacket)
	r.FractionLost = rawPacket[fractionLostOffset]

	tlBytes := rawPacket[totalLostOffset:]
	uTotalLost := uint32(tlBytes[2]) | uint32(tlBytes[1])<<8 | uint32(tlBytes[0])<<16
	// test sign
	ua := uTotalLost & 0x007fffff
	if uTotalLost&0x00800000 == 0x00800000 {
		ua |= 0xff800000
	}
	r.TotalLost = int32(ua)

	r.LastSequenceNumber = binary.BigEndian.Uint32(rawPacket[lastSeqOffset:])
	r.Jitter = binary.BigEndian.Uint32(rawPacket[jitterOffset:])
	r.LastSenderReport = binary.BigEndian.Uint32(rawPacket[lastSROffset:])
	r.Delay = binary.BigEndian.Uint32(rawPacket[delayOffset:])

	return nil
}

func (r *ReceptionReport) len() int {
	return receptionReportLength
}
