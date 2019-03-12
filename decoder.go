package rtcp

import "io"

// A Decoder reads packets from an RTCP combined packet.
type Decoder struct {
	r io.Reader
}

// NewDecoder creates a new Decoder reading from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

// DecodePacket reads one packet from r.
//
// It returns the parsed packet Header and a byte slice containing the encoded
// packet data (including the header). How the packet data is parsed depends on
// the Type field contained in the Header.
func (r *Decoder) DecodePacket() (Packet, error) {
	// First grab the header
	var header Header
	headerBuf := make([]byte, headerLength)
	if _, err := io.ReadFull(r.r, headerBuf); err != nil {
		return nil, err
	}
	if err := header.Unmarshal(headerBuf); err != nil {
		return nil, err
	}

	packetLen := (header.Length + 1) * 4

	// Then grab the rest
	bodyBuf := make([]byte, packetLen-headerLength)
	if _, err := io.ReadFull(r.r, bodyBuf); err != nil {
		return nil, err
	}
	data := append(headerBuf, bodyBuf...)

	p, bytesused, err := unmarshal(data)

	//data should be exactly one RTCP packet.
	if bytesused != len(data) {
		return nil, errInvalidHeader
	}

	if err != nil {
		return nil, err
	}

	return p, err
}
