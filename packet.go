package rtcp

// Packet represents an RTCP packet, a protocol used for out-of-band statistics and control information for an RTP session
type Packet interface {
	Header() Header
	// DestinationSSRC returns an array of SSRC values that this packet refers to.
	DestinationSSRC() []uint32

	Marshal() ([]byte, error)
	Unmarshal(rawPacket []byte) error
}

// A CompoundPacket is a collection of RTCP packets transmitted as a single packet with
// the underlying protocol (for example UDP).
//
// To maximize the resolution of receiption statistics, the first Packet in a CompoundPacket
// must always be either a SenderReport or a ReceiverReport.  This is true even if no data
// has been sent or received, in which case an empty ReceiverReport must be sent, and even
// if the only other RTCP packet in the compound packet is a Goodbye.
//
// Next, a SourceDescription containing a CNAME item must be included in each CompoundPacket
// to identify the source and to begin associating media for purposes such as lip-sync.
//
// Other RTCP packet types may follow in any order. Packet types may appear more than once.
type CompoundPacket []Packet

// unmarshal is a factory which pulls the first RTCP packet from a bytestream,
// and returns it's parsed representation, and the amount of data that was processed.
func unmarshal(rawData []byte) (packet Packet, bytesprocessed int, err error) {
	var h Header

	err = h.Unmarshal(rawData)
	if err != nil {
		return nil, 0, err
	}

	bytesprocessed = int(h.Length+1) * 4
	inPacket := rawData[:bytesprocessed]

	switch h.Type {
	case TypeSenderReport:
		packet = new(SenderReport)

	case TypeReceiverReport:
		packet = new(ReceiverReport)

	case TypeSourceDescription:
		packet = new(SourceDescription)

	case TypeGoodbye:
		packet = new(Goodbye)

	case TypeTransportSpecificFeedback:
		switch h.Count {
		case FormatTLN:
			packet = new(TransportLayerNack)
		case FormatRRR:
			packet = new(RapidResynchronizationRequest)
		default:
			packet = new(RawPacket)
		}

	case TypePayloadSpecificFeedback:
		switch h.Count {
		case FormatPLI:
			packet = new(PictureLossIndication)
		case FormatSLI:
			packet = new(SliceLossIndication)
		case FormatREMB:
			packet = new(ReceiverEstimatedMaximumBitrate)
		default:
			packet = new(RawPacket)
		}

	default:
		packet = new(RawPacket)
	}

	err = packet.Unmarshal(inPacket)

	return packet, bytesprocessed, err
}

// Unmarshal takes an entire udp datagram (which may consist of multiple RTCP packets) and returns
// an unmarshalled array of packets.
func Unmarshal(rawData []byte) (CompoundPacket, error) {
	var out CompoundPacket

	for len(rawData) != 0 {
		p, processed, err := unmarshal(rawData)

		if err != nil {
			return nil, err
		}

		out = append(out, p)
		rawData = rawData[processed:]
	}

	var err error

	//some extra validity checks for compound packets
	//(if they fail, return the (now successfully parsed) packets, but an error too)
	if len(out) > 1 {
		if out[0].Header().Padding {
			//padding isn't allowed in the first packet in a compound datagram
			err = errInvalidHeader
		} else if (out[0].Header().Type != TypeSenderReport) &&
			(out[0].Header().Type != TypeReceiverReport) {
			//SenderReport and ReceiverReport are the only types that
			//are allowed to be the first packet in a compound datagram
			err = errInvalidHeader
		}
	}

	return out, err
}
