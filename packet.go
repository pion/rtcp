// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"bytes"
	"sync"
)

// Packet represents an RTCP packet, a protocol used for out-of-band statistics and control information for an RTP session
type Packet interface {
	// DestinationSSRC returns an array of SSRC values that this packet refers to.
	DestinationSSRC() []uint32

	Marshal() ([]byte, error)
	Unmarshal(rawPacket []byte) error
	MarshalSize() int

	// Release returns the packet to its pool
	Release()
}

var (
	senderReportPool                    = sync.Pool{New: func() interface{} { return new(SenderReport) }}
	receiverReportPool                  = sync.Pool{New: func() interface{} { return new(ReceiverReport) }}
	sourceDescriptionPool               = sync.Pool{New: func() interface{} { return new(SourceDescription) }}
	goodbyePool                         = sync.Pool{New: func() interface{} { return new(Goodbye) }}
	transportLayerNackPool              = sync.Pool{New: func() interface{} { return new(TransportLayerNack) }}
	rapidResynchronizationRequestPool   = sync.Pool{New: func() interface{} { return new(RapidResynchronizationRequest) }}
	transportLayerCCPool                = sync.Pool{New: func() interface{} { return new(TransportLayerCC) }}
	ccFeedbackReportPool                = sync.Pool{New: func() interface{} { return new(CCFeedbackReport) }}
	pictureLossIndicationPool           = sync.Pool{New: func() interface{} { return new(PictureLossIndication) }}
	sliceLossIndicationPool             = sync.Pool{New: func() interface{} { return new(SliceLossIndication) }}
	receiverEstimatedMaximumBitratePool = sync.Pool{New: func() interface{} { return new(ReceiverEstimatedMaximumBitrate) }}
	fullIntraRequestPool                = sync.Pool{New: func() interface{} { return new(FullIntraRequest) }}
	extendedReportPool                  = sync.Pool{New: func() interface{} { return new(ExtendedReport) }}
	applicationDefinedPool              = sync.Pool{New: func() interface{} { return new(ApplicationDefined) }}
	rawPacketPool                       = sync.Pool{New: func() interface{} { return new(RawPacket) }}
)

// Unmarshal takes an entire udp datagram (which may consist of multiple RTCP packets) and
// returns the unmarshaled packets it contains.
//
// If this is a reduced-size RTCP packet a feedback packet (Goodbye, SliceLossIndication, etc)
// will be returned. Otherwise, the underlying type of the returned packet will be
// CompoundPacket.
func Unmarshal(rawData []byte) ([]Packet, error) {
	estimatedPackets := len(rawData) / 100 // Estimate based on average packet size
	packets := make([]Packet, 0, estimatedPackets)
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

// Marshal takes an array of Packets and serializes them to a single buffer
func Marshal(packets []Packet) ([]byte, error) {
	var buf bytes.Buffer
	for _, p := range packets {
		data, err := p.Marshal()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		p.Release()
	}
	return buf.Bytes(), nil
}

// unmarshal is a factory which pulls the first RTCP packet from a bytestream,
// and returns it's parsed representation, and the amount of data that was processed.
func unmarshal(rawData []byte) (packet Packet, bytesprocessed int, err error) {
	var h Header

	err = h.Unmarshal(rawData)
	if err != nil {
		return nil, 0, err
	}

	bytesprocessed = int(h.Length+1) * 4
	if bytesprocessed > len(rawData) {
		return nil, 0, errPacketTooShort
	}
	inPacket := rawData[:bytesprocessed]

	switch h.Type {
	case TypeSenderReport:
		packet = senderReportPool.Get().(*SenderReport)

	case TypeReceiverReport:
		packet = receiverReportPool.Get().(*ReceiverReport)

	case TypeSourceDescription:
		packet = sourceDescriptionPool.Get().(*SourceDescription)

	case TypeGoodbye:
		packet = goodbyePool.Get().(*Goodbye)

	case TypeTransportSpecificFeedback:
		switch h.Count {
		case FormatTLN:
			packet = transportLayerNackPool.Get().(*TransportLayerNack)
		case FormatRRR:
			packet = rapidResynchronizationRequestPool.Get().(*RapidResynchronizationRequest)
		case FormatTCC:
			packet = transportLayerCCPool.Get().(*TransportLayerCC)
		case FormatCCFB:
			packet = ccFeedbackReportPool.Get().(*CCFeedbackReport)
		default:
			packet = rawPacketPool.Get().(*RawPacket)
		}

	case TypePayloadSpecificFeedback:
		switch h.Count {
		case FormatPLI:
			packet = pictureLossIndicationPool.Get().(*PictureLossIndication)
		case FormatSLI:
			packet = sliceLossIndicationPool.Get().(*SliceLossIndication)
		case FormatREMB:
			packet = receiverEstimatedMaximumBitratePool.Get().(*ReceiverEstimatedMaximumBitrate)
		case FormatFIR:
			packet = fullIntraRequestPool.Get().(*FullIntraRequest)
		default:
			packet = rawPacketPool.Get().(*RawPacket)
		}

	case TypeExtendedReport:
		packet = extendedReportPool.Get().(*ExtendedReport)

	case TypeApplicationDefined:
		packet = applicationDefinedPool.Get().(*ApplicationDefined)

	default:
		packet = rawPacketPool.Get().(*RawPacket)
	}

	err = packet.Unmarshal(inPacket)

	return packet, bytesprocessed, err
}

// Implement the Release method for each concrete packet type
func (p *SenderReport) Release() {
	senderReportPool.Put(p)
}

func (p *ReceiverReport) Release() {
	receiverReportPool.Put(p)
}

func (p *SourceDescription) Release() {
	sourceDescriptionPool.Put(p)
}

func (p *Goodbye) Release() {
	goodbyePool.Put(p)
}

func (p *TransportLayerNack) Release() {
	transportLayerNackPool.Put(p)
}

func (p *RapidResynchronizationRequest) Release() {
	rapidResynchronizationRequestPool.Put(p)
}

func (p *TransportLayerCC) Release() {
	transportLayerCCPool.Put(p)
}

func (p *CCFeedbackReport) Release() {
	ccFeedbackReportPool.Put(p)
}

func (p *PictureLossIndication) Release() {
	pictureLossIndicationPool.Put(p)
}

func (p *SliceLossIndication) Release() {
	sliceLossIndicationPool.Put(p)
}

func (p *ReceiverEstimatedMaximumBitrate) Release() {
	receiverEstimatedMaximumBitratePool.Put(p)
}

func (p *FullIntraRequest) Release() {
	fullIntraRequestPool.Put(p)
}

func (p *ExtendedReport) Release() {
	extendedReportPool.Put(p)
}

func (p *ApplicationDefined) Release() {
	applicationDefinedPool.Put(p)
}

func (p *CompoundPacket) Release() {
	// CompoundPacket is a slice of pointers, so we need to release each one
	for _, packet := range *p {
		packet.Release()
	}
}

func (p *RawPacket) Release() {
	rawPacketPool.Put(p)
}
