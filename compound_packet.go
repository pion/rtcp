package rtcp

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

	// some extra validity checks for compound packets
	// (if they fail, return the (now successfully parsed) packets, but an error too)
	if len(out) == 0 {
		return out, errInvalidHeader
	}

	firstHdr := out[0].Header()
	if firstHdr.Padding {
		// padding isn't allowed in the first packet in a compound datagram
		return out, errInvalidHeader
	} else if (firstHdr.Type != TypeSenderReport) &&
		(firstHdr.Type != TypeReceiverReport) {
		// SenderReport and ReceiverReport are the only types that
		// are allowed to be the first packet in a compound datagram
		return out, errInvalidHeader
	}

	return out, err
}
