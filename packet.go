package rtcp

// Packet represents an RTCP packet, a protocol used for out-of-band statistics and control information for an RTP session
type Packet interface {
	Header() Header
	// DestinationSSRC returns an array of SSRC values that this packet refers to.
	DestinationSSRC() []uint32

	Marshal() ([]byte, error)
	Unmarshal(rawPacket []byte) error
}

// Unmarshal is a factory a polymorphic RTCP packet, and its header,
func Unmarshal(rawData []byte) ([]Packet, error) {
	var err error
	var out []Packet

	for len(rawData) != 0 {
		var h Header
		var p Packet

		err = h.Unmarshal(rawData)
		if err != nil {
			return nil, err
		}

		plen := (h.Length + 1) * 4
		inPacket := rawData[:plen]
		rawData = rawData[plen:]

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
		out = append(out, p)
	}
	return out, err
}
