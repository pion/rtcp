// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// An RTCP packet from a packet dump
func realPacket() []byte {
	return []byte{
		// Receiver Report (offset=0)
		// v=2, p=0, count=1, RR, len=7
		0x81, 0xc9, 0x0, 0x7,
		// ssrc=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,
		// ssrc=0xbc5e9a40
		0xbc, 0x5e, 0x9a, 0x40,
		// fracLost=0, totalLost=0
		0x0, 0x0, 0x0, 0x0,
		// lastSeq=0x46e1
		0x0, 0x0, 0x46, 0xe1,
		// jitter=273
		0x0, 0x0, 0x1, 0x11,
		// lsr=0x9f36432
		0x9, 0xf3, 0x64, 0x32,
		// delay=150137
		0x0, 0x2, 0x4a, 0x79,

		// Source Description (offset=32)
		// v=2, p=0, count=1, SDES, len=12
		0x81, 0xca, 0x0, 0xc,
		// ssrc=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,
		// CNAME, len=38
		0x1, 0x26,
		// text="{9c00eb92-1afb-9d49-a47d-91f64eee69f5}"
		0x7b, 0x39, 0x63, 0x30,
		0x30, 0x65, 0x62, 0x39,
		0x32, 0x2d, 0x31, 0x61,
		0x66, 0x62, 0x2d, 0x39,
		0x64, 0x34, 0x39, 0x2d,
		0x61, 0x34, 0x37, 0x64,
		0x2d, 0x39, 0x31, 0x66,
		0x36, 0x34, 0x65, 0x65,
		0x65, 0x36, 0x39, 0x66,
		0x35, 0x7d,
		// END + padding
		0x0, 0x0, 0x0, 0x0,

		// Goodbye (offset=84)
		// v=2, p=0, count=1, BYE, len=1
		0x81, 0xcb, 0x0, 0x1,
		// source=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,

		// Picture Loss Indication (offset=92)
		0x81, 0xce, 0x0, 0x2,
		// sender=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,
		// media=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,

		// RapidResynchronizationRequest (offset=104)
		0x85, 0xcd, 0x0, 0x2,
		// sender=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,
		// media=0x902f9e2e
		0x90, 0x2f, 0x9e, 0x2e,

		// ApplicationDefined (offset=116)
		0x80, 0xcc, 0x00, 0x03,
		// sender=0x4baae1ab
		0x4b, 0xaa, 0xe1, 0xab,
		// name='NAME'
		0x4E, 0x41, 0x4D, 0x45,
		// data='ABCD'
		0x41, 0x42, 0x43, 0x44,
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	packetData := realPacket()
	for i := 0; i < b.N; i++ {
		pkts, err := Unmarshal(packetData)
		if err != nil {
			b.Fatalf("Error unmarshalling packets: %s", err)
		}

		for _, pkt := range pkts {
			pkt.Release()
		}

	}
}

func TestUnmarshal(t *testing.T) {
	packet, err := Unmarshal(realPacket())
	if err != nil {
		t.Fatalf("Error unmarshalling packets: %s", err)
	}

	expected := []Packet{
		&ReceiverReport{
			SSRC: 0x902f9e2e,
			Reports: []ReceptionReport{{
				SSRC:               0xbc5e9a40,
				FractionLost:       0,
				TotalLost:          0,
				LastSequenceNumber: 0x46e1,
				Jitter:             273,
				LastSenderReport:   0x9f36432,
				Delay:              150137,
			}},
			ProfileExtensions: []byte{},
		},
		NewCNAMESourceDescription(0x902f9e2e, "{9c00eb92-1afb-9d49-a47d-91f64eee69f5}"),
		&Goodbye{
			Sources: []uint32{0x902f9e2e},
		},
		&PictureLossIndication{
			SenderSSRC: 0x902f9e2e,
			MediaSSRC:  0x902f9e2e,
		},
		&RapidResynchronizationRequest{
			SenderSSRC: 0x902f9e2e,
			MediaSSRC:  0x902f9e2e,
		},
		&ApplicationDefined{
			SSRC: 0x4baae1ab,
			Name: "NAME",
			Data: []byte{0x41, 0x42, 0x43, 0x44},
		},
	}

	assert.Equal(t, expected, packet)
}

func TestUnmarshalNil(t *testing.T) {
	_, err := Unmarshal(nil)
	if got, want := err, errInvalidHeader; !errors.Is(got, want) {
		t.Fatalf("Unmarshal(nil) err = %v, want %v", got, want)
	}
}

func TestInvalidHeaderLength(t *testing.T) {
	invalidPacket := []byte{
		// Receiver Report (offset=0)
		// v=2, p=0, count=1, RR, len=100
		0x81, 0xc9, 0x0, 0x64,
	}

	_, err := Unmarshal(invalidPacket)
	if got, want := err, errPacketTooShort; !errors.Is(got, want) {
		t.Fatalf("Unmarshal(nil) err = %v, want %v", got, want)
	}
}

func TestPacketPool(t *testing.T) {
	t.Run("SenderReport", func(t *testing.T) {
		sr := senderReportPool.Get()
		p, ok := sr.(*SenderReport)
		assert.True(t, ok)

		p.Release()
	})

	t.Run("ReceiverReport", func(t *testing.T) {
		rr := receiverReportPool.Get()
		p, ok := rr.(*ReceiverReport)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("SourceDescription", func(t *testing.T) {
		sd := sourceDescriptionPool.Get()
		p, ok := sd.(*SourceDescription)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("Goodbye", func(t *testing.T) {
		gb := goodbyePool.Get()
		p, ok := gb.(*Goodbye)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("TransportLayerNack", func(t *testing.T) {
		tln := transportLayerNackPool.Get()
		p, ok := tln.(*TransportLayerNack)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("RapidResynchronizationRequest", func(t *testing.T) {
		rrr := rapidResynchronizationRequestPool.Get()
		p, ok := rrr.(*RapidResynchronizationRequest)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("TransportLayerCC", func(t *testing.T) {
		tcc := transportLayerCCPool.Get()
		p, ok := tcc.(*TransportLayerCC)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("CCFeedbackReport", func(t *testing.T) {
		ccfb := ccFeedbackReportPool.Get()
		p, ok := ccfb.(*CCFeedbackReport)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("PictureLossIndication", func(t *testing.T) {
		pli := pictureLossIndicationPool.Get()
		p, ok := pli.(*PictureLossIndication)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("SliceLossIndication", func(t *testing.T) {
		sli := sliceLossIndicationPool.Get()
		p, ok := sli.(*SliceLossIndication)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("ReceiverEstimatedMaximumBitrate", func(t *testing.T) {
		remb := receiverEstimatedMaximumBitratePool.Get()
		p, ok := remb.(*ReceiverEstimatedMaximumBitrate)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("FullIntraRequest", func(t *testing.T) {
		fir := fullIntraRequestPool.Get()
		p, ok := fir.(*FullIntraRequest)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("ExtendedReport", func(t *testing.T) {
		er := extendedReportPool.Get()
		p, ok := er.(*ExtendedReport)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("ApplicationDefined", func(t *testing.T) {
		ad := applicationDefinedPool.Get()
		p, ok := ad.(*ApplicationDefined)
		assert.True(t, ok)
		p.Release()
	})

	t.Run("RawPacket", func(t *testing.T) {
		rp := rawPacketPool.Get()
		p, ok := rp.(*RawPacket)
		assert.True(t, ok)
		p.Release()
	})
}
