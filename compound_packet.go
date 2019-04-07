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

var _ Packet = (*CompoundPacket)(nil) // assert is a Packet

// Validate returns an error if this is not an RFC-compliant CompoundPacket.
func (c CompoundPacket) Validate() error {
	if len(c) == 0 {
		return errEmptyCompound
	}

	// SenderReport and ReceiverReport are the only types that
	// are allowed to be the first packet in a compound datagram
	switch c[0].(type) {
	case *SenderReport, *ReceiverReport:
		// ok
	default:
		return errBadFirstPacket
	}

	for _, pkt := range c[1:] {
		switch p := pkt.(type) {
		// If the number of RecetpionReports exceeds 31 additional ReceiverReports
		// can be included here.
		case *ReceiverReport:
			continue

		// A SourceDescription containing a CNAME must be included in every
		// CompoundPacket.
		case *SourceDescription:
			var hasCNAME bool
			for _, c := range p.Chunks {
				for _, it := range c.Items {
					if it.Type == SDESCNAME {
						hasCNAME = true
					}
				}
			}

			if !hasCNAME {
				return errMissingCNAME
			}

			return nil

		// Other packets are not permitted before the CNAME
		default:
			return errPacketBeforeCNAME
		}
	}

	// CNAME never reached
	return errMissingCNAME
}

// Marshal encodes the CompoundPacket as binary.
func (c CompoundPacket) Marshal() ([]byte, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	out := make([]byte, 0)
	for _, p := range c {
		data, err := p.Marshal()
		if err != nil {
			return nil, err
		}
		out = append(out, data...)
	}
	return out, nil
}

// Unmarshal decodes a CompoundPacket from binary.
func (c *CompoundPacket) Unmarshal(rawData []byte) error {
	out := make(CompoundPacket, 0)
	for len(rawData) != 0 {
		p, processed, err := unmarshal(rawData)

		if err != nil {
			return err
		}

		out = append(out, p)
		rawData = rawData[processed:]
	}
	*c = out

	if err := c.Validate(); err != nil {
		return err
	}

	return nil
}

// DestinationSSRC returns the synchronization sources associated with this
// CompoundPacket's reception report.
func (c CompoundPacket) DestinationSSRC() []uint32 {
	if len(c) == 0 {
		return nil
	}

	return c[0].DestinationSSRC()
}
