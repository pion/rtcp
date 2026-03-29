// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import "fmt"

var packetTypeNames = map[PacketType]string{
	TypeSenderReport:       "SenderReport",
	TypeReceiverReport:     "ReceiverReport",
	TypeSourceDescription:  "SourceDescription",
	TypeGoodbye:            "Goodbye",
	TypeExtendedReport:     "ExtendedReport",
	TypeApplicationDefined: "ApplicationDefined",
}

var transportSpecificFeedbackNames = map[uint8]string{
	FormatTLN:  "TransportLayerNack",
	FormatRRR:  "RapidResynchronizationRequest",
	FormatTCC:  "TransportLayerCC",
	FormatCCFB: "CCFeedbackReport",
}

var payloadSpecificFeedbackNames = map[uint8]string{
	FormatPLI:  "PictureLossIndication",
	FormatSLI:  "SliceLossIndication",
	FormatREMB: "ReceiverEstimatedMaximumBitrate",
	FormatFIR:  "FullIntraRequest",
}

// Packet represents an RTCP packet, a protocol used for out-of-band statistics
// and control information for an RTP session.
type Packet interface {
	// DestinationSSRC returns an array of SSRC values that this packet refers to.
	DestinationSSRC() []uint32

	Marshal() ([]byte, error)
	Unmarshal(rawPacket []byte) error
	MarshalSize() int
}

// Unmarshal takes an entire udp datagram (which may consist of multiple RTCP packets) and
// returns the unmarshaled packets it contains.
//
// If this is a reduced-size RTCP packet a feedback packet (Goodbye, SliceLossIndication, etc)
// will be returned. Otherwise, the underlying type of the returned packet will be
// CompoundPacket.
func Unmarshal(rawData []byte) ([]Packet, error) {
	var packets []Packet
	for len(rawData) != 0 {
		p, processed, err := unmarshal(rawData)
		if err != nil {
			return nil, err
		}

		packets = append(packets, p)
		rawData = rawData[processed:]
	}

	switch len(packets) {
	// Empty packet
	case 0:
		return nil, errInvalidHeader
	// Multiple Packets
	default:
		return packets, nil
	}
}

// Marshal takes an array of Packets and serializes them to a single buffer.
func Marshal(packets []Packet) ([]byte, error) {
	out := make([]byte, 0)
	for _, p := range packets {
		data, err := p.Marshal()
		if err != nil {
			return nil, err
		}
		out = append(out, data...)
	}

	return out, nil
}

// unmarshal is a factory which pulls the first RTCP packet from a bytestream,
// and returns it's parsed representation, and the amount of data that was processed.
//
//nolint:cyclop
func unmarshal(rawData []byte) (packet Packet, bytesprocessed int, err error) {
	var header Header

	err = header.Unmarshal(rawData)
	if err != nil {
		return nil, 0, err
	}

	bytesprocessed = int(header.Length+1) * 4
	if bytesprocessed > len(rawData) {
		return nil, 0, errPacketTooShortFor(packetNameFromHeader(header))
	}
	inPacket := rawData[:bytesprocessed]

	switch header.Type {
	case TypeSenderReport:
		packet = new(SenderReport)

	case TypeReceiverReport:
		packet = new(ReceiverReport)

	case TypeSourceDescription:
		packet = new(SourceDescription)

	case TypeGoodbye:
		packet = new(Goodbye)

	case TypeTransportSpecificFeedback:
		switch header.Count {
		case FormatTLN:
			packet = new(TransportLayerNack)
		case FormatRRR:
			packet = new(RapidResynchronizationRequest)
		case FormatTCC:
			packet = new(TransportLayerCC)
		case FormatCCFB:
			packet = new(CCFeedbackReport)
		default:
			packet = new(RawPacket)
		}

	case TypePayloadSpecificFeedback:
		switch header.Count {
		case FormatPLI:
			packet = new(PictureLossIndication)
		case FormatSLI:
			packet = new(SliceLossIndication)
		case FormatREMB:
			packet = new(ReceiverEstimatedMaximumBitrate)
		case FormatFIR:
			packet = new(FullIntraRequest)
		default:
			packet = new(RawPacket)
		}

	case TypeExtendedReport:
		packet = new(ExtendedReport)

	case TypeApplicationDefined:
		packet = new(ApplicationDefined)

	default:
		packet = new(RawPacket)
	}

	err = packet.Unmarshal(inPacket)

	return packet, bytesprocessed, err
}

func packetNameFromHeader(header Header) string {
	if header.Type == TypeTransportSpecificFeedback {
		return transportSpecificFeedbackName(header.Count)
	}

	if header.Type == TypePayloadSpecificFeedback {
		return payloadSpecificFeedbackName(header.Count)
	}

	if name, ok := packetTypeNames[header.Type]; ok {
		return name
	}

	return fmt.Sprintf("PacketType(%d)", header.Type)
}

func transportSpecificFeedbackName(count uint8) string {
	if name, ok := transportSpecificFeedbackNames[count]; ok {
		return name
	}

	return fmt.Sprintf("TransportSpecificFeedback(FMT=%d)", count)
}

func payloadSpecificFeedbackName(count uint8) string {
	if name, ok := payloadSpecificFeedbackNames[count]; ok {
		return name
	}

	return fmt.Sprintf("PayloadSpecificFeedback(FMT=%d)", count)
}
