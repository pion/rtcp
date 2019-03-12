package rtcp

// Packet represents an RTCP packet, a protocol used for out-of-band statistics and control information for an RTP session
type Packet interface {
	Header() Header
	// DestinationSSRC returns an array of SSRC values that this packet refers to.
	DestinationSSRC() []uint32

	Marshal() ([]byte, error)
	Unmarshal(rawPacket []byte) error
}

// Unmarshal is a factory which pulls the first RTCP packet from a bytestream,
// and returns it's parsed representation, and the amount of data that was processed.
func Unmarshal(rawData []byte) (Packet, int, error) {
	var err error
	var h Header
	var p Packet

	err = h.Unmarshal(rawData)
	if err != nil {
		return nil, 0, err
	}

	plen := (h.Length + 1) * 4
	inPacket := rawData[:plen]

	switch h.Type {
	case TypeSenderReport:
		p = new(SenderReport)

	case TypeReceiverReport:
		p = new(ReceiverReport)

	case TypeSourceDescription:
		p = new(SourceDescription)

	case TypeGoodbye:
		p = new(Goodbye)

	case TypeTransportSpecificFeedback:
		switch h.Count {
		case FormatTLN:
			p = new(TransportLayerNack)
		case FormatRRR:
			p = new(RapidResynchronizationRequest)
		default:
			p = new(RawPacket)
		}

	case TypePayloadSpecificFeedback:
		switch h.Count {
		case FormatPLI:
			p = new(PictureLossIndication)
		case FormatSLI:
			p = new(SliceLossIndication)
		case FormatREMB:
			p = new(ReceiverEstimatedMaximumBitrate)
		default:
			p = new(RawPacket)
		}

	default:
		p = new(RawPacket)
	}

	err = p.Unmarshal(inPacket)

	return p, int(plen), err
}

// UnmarshalDatagram takes an entire udp datagram (which may consist of multiple RTCP packets) and returns
// an unmarshalled array of packets.
func UnmarshalDatagram(rawData []byte) ([]Packet, error) {
	var out []Packet

	for len(rawData) != 0 {
		p, processed, err := Unmarshal(rawData)
		if err != nil {
			return nil, err
		}

		out = append(out, p)
		rawData = rawData[processed:]
	}
	return out, nil
}
