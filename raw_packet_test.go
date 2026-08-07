// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Packet = (*RawPacket)(nil) // assert is a Packet

func TestRawPacketRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name               string
		Packet             RawPacket
		WantMarshalError   error
		WantUnmarshalError error
	}{
		{
			Name: "valid",
			Packet: RawPacket([]byte{
				// v=2, p=0, count=1, BYE, len=12
				0x81, 0xcb, 0x00, 0x0c,
				// ssrc=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// len=3, text=FOO
				0x03, 0x46, 0x4f, 0x4f,
			}),
		},
		{
			Name:               "short header",
			Packet:             RawPacket([]byte{0x00}),
			WantUnmarshalError: errPacketTooShort,
		},
		{
			Name: "invalid header",
			Packet: RawPacket([]byte{
				// v=0, p=0, count=0, RR, len=4
				0x00, 0xc9, 0x00, 0x04,
			}),
			WantUnmarshalError: errBadVersion,
		},
	} {
		data, err := test.Packet.Marshal()
		assert.ErrorIsf(t, err, test.WantMarshalError, "Marshal %q", test.Name)
		if err != nil {
			continue
		}

		var decoded RawPacket

		err = decoded.Unmarshal(data)
		assert.ErrorIsf(t, err, test.WantUnmarshalError, "Unmarshal %q", test.Name)
		if err != nil {
			continue
		}
		assert.Equalf(t, test.Packet, decoded, "Unmarshal %q", test.Name)
	}
}

func destinationSSRCCorpus(tb testing.TB) [][]byte {
	tb.Helper()

	packets := [][]Packet{
		{&SenderReport{SSRC: 0xAAAA0001, Reports: []ReceptionReport{{SSRC: 0xBBBB0001}, {SSRC: 0xBBBB0002}}}},
		{&ReceiverReport{SSRC: 0xAAAA0002, Reports: []ReceptionReport{{SSRC: 0xBBBB0001}, {SSRC: 0xBBBB0003}}}},
		{&SourceDescription{Chunks: []SourceDescriptionChunk{
			{Source: 0xCCCC0001, Items: []SourceDescriptionItem{{Type: SDESCNAME, Text: "alice@example"}}},
			{Source: 0xCCCC0002, Items: []SourceDescriptionItem{{Type: SDESCNAME, Text: "bob"}}},
		}}},
		{&Goodbye{Sources: []uint32{0xDDDD0001, 0xDDDD0002}, Reason: "shutdown"}},
		{&ApplicationDefined{SSRC: 0xEEEE0001, Name: "TEST", Data: []byte{1, 2, 3, 4}}},
		{&PictureLossIndication{SenderSSRC: 0xAAAA0001, MediaSSRC: 0xBBBB0001}},
		{&SliceLossIndication{
			SenderSSRC: 0xAAAA0001,
			MediaSSRC:  0xBBBB0002,
			SLI:        []SLIEntry{{First: 1, Number: 2, Picture: 3}},
		}},
		{&FullIntraRequest{
			SenderSSRC: 0xAAAA0001,
			MediaSSRC:  0xBBBB0001,
			FIR:        []FIREntry{{SSRC: 0xBBBB0001, SequenceNumber: 1}, {SSRC: 0xBBBB0002, SequenceNumber: 2}},
		}},
		{&ReceiverEstimatedMaximumBitrate{SenderSSRC: 0xAAAA0001, Bitrate: 1e6, SSRCs: []uint32{0xBBBB0001, 0xBBBB0002}}},
		{&ReceiverEstimatedMaximumBitrate{SenderSSRC: 0xAAAA0001, Bitrate: 1e6}},
		{&TransportLayerNack{
			SenderSSRC: 0xAAAA0001,
			MediaSSRC:  0xBBBB0001,
			Nacks:      []NackPair{{PacketID: 42, LostPackets: 2}},
		}},
		{&RapidResynchronizationRequest{SenderSSRC: 0xAAAA0001, MediaSSRC: 0xBBBB0001}},
		{&CCFeedbackReport{SenderSSRC: 0xAAAA0001, ReportBlocks: []CCFeedbackReportBlock{
			{MediaSSRC: 0xBBBB0001, BeginSequence: 1, MetricBlocks: []CCFeedbackMetricBlock{
				{Received: true, ArrivalTimeOffset: 12}, {Received: true, ArrivalTimeOffset: 34}, {Received: false},
			}},
			{MediaSSRC: 0xBBBB0002, BeginSequence: 7, MetricBlocks: []CCFeedbackMetricBlock{
				{Received: true, ArrivalTimeOffset: 56},
			}},
		}, ReportTimestamp: 0x01020304}},
		{&ExtendedReport{SenderSSRC: 0xAAAA0001, Reports: []ReportBlock{
			// An even number of chunks keeps the block 32-bit aligned.
			&LossRLEReportBlock{SSRC: 0xBBBB0001, Chunks: []Chunk{0x4006, 0x0006}},
			&DLRRReportBlock{Reports: []DLRRReport{
				{SSRC: 0xBBBB0002, LastRR: 1, DLRR: 2}, {SSRC: 0xBBBB0003, LastRR: 3, DLRR: 4},
			}},
			&ReceiverReferenceTimeReportBlock{NTPTimestamp: 0x0102030405060708},
			&DuplicateRLEReportBlock{SSRC: 0xBBBB0004, Chunks: []Chunk{0x4006, 0x0006}},
			&PacketReceiptTimesReportBlock{SSRC: 0xBBBB0005},
			&StatisticsSummaryReportBlock{SSRC: 0xBBBB0006},
			&VoIPMetricsReportBlock{SSRC: 0xBBBB0007},
		}}},
	}

	var datagrams [][]byte
	for _, pkts := range packets {
		raw, err := Marshal(pkts)
		assert.NoErrorf(tb, err, "marshal corpus packets %T", pkts[0])
		datagrams = append(datagrams, raw)
	}

	// RR + TWCC, the typical receiver->sender feedback compound. The TWCC
	// packet carries one run-length chunk of two received small deltas.
	twcc := []byte{
		0x8f, 0xcd, 0x00, 0x05, // V=2, FMT=15, PT=205, length=5
		0x11, 0x11, 0x11, 0x11, // sender SSRC
		0x22, 0x22, 0x22, 0x22, // media SSRC
		0x03, 0xe8, 0x00, 0x02, // base sequence 1000, packet status count 2
		0x01, 0x23, 0x45, 0x03, // reference time, fb pkt count
		0x20, 0x02, 0x04, 0x08, // run-length chunk, recv deltas 1ms and 2ms
	}
	rrAndTWCC := append(append([]byte{}, datagrams[1]...), twcc...)
	datagrams = append(datagrams, rrAndTWCC)

	// An unknown packet type (195), parsed as a RawPacket with no
	// destinations, alone and inside a compound.
	unknown := []byte{0x80, 195, 0x00, 0x01, 0xDE, 0xAD, 0xBE, 0xEF}
	datagrams = append(datagrams, unknown)

	// An unknown feedback format (2) within a known packet type (205), also
	// parsed as a RawPacket.
	datagrams = append(datagrams, []byte{0x82, 0xCD, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
	datagrams = append(datagrams, append(append([]byte{}, datagrams[0]...), unknown...))

	datagrams = append(datagrams, realPacket())

	// A datagram Unmarshal rejects would be silently skipped by the tests
	// comparing against it.
	for i, datagram := range datagrams {
		_, err := Unmarshal(datagram)
		assert.NoErrorf(tb, err, "corpus datagram %d does not unmarshal", i)
	}

	return datagrams
}

// ParseDestinationSSRC must report the same destination SSRCs per packet as
// a full unmarshal does.
func TestParseDestinationSSRCMatchesUnmarshal(t *testing.T) {
	for di, datagram := range destinationSSRCCorpus(t) {
		raws, err := UnmarshalRaw(datagram)
		assert.NoErrorf(t, err, "datagram %d", di)

		pkts, err := Unmarshal(datagram)
		assert.NoErrorf(t, err, "datagram %d", di)
		if !assert.Len(t, raws, len(pkts)) {
			continue
		}

		for i, pkt := range pkts {
			got, err := raws[i].ParseDestinationSSRC(nil)
			assert.NoErrorf(t, err, "datagram %d packet %d (%T)", di, i, pkt)
			assert.ElementsMatchf(t, pkt.DestinationSSRC(), got, "destination SSRCs of datagram %d packet %d (%T)", di, i, pkt)
		}
	}
}

// Every packet type the Unmarshal factory can produce, RawPacket included,
// must appear in the corpus, so TestParseDestinationSSRCMatchesUnmarshal
// cannot silently lose coverage when a new packet type or format is added.
func TestRawPacketCorpusCoverage(t *testing.T) {
	corpusTypes := map[reflect.Type]bool{}
	for _, datagram := range destinationSSRCCorpus(t) {
		pkts, err := Unmarshal(datagram)
		assert.NoError(t, err)
		for _, pkt := range pkts {
			corpusTypes[reflect.TypeOf(pkt)] = true
		}
	}

	for packetType := 0; packetType <= 0xFF; packetType++ {
		formats := 1
		if PacketType(packetType) == TypeTransportSpecificFeedback || PacketType(packetType) == TypePayloadSpecificFeedback {
			formats = 32
		}
		for format := range formats {
			// The factory picks the concrete type before unmarshaling, so a
			// bare header is enough to learn the mapping.
			//nolint:gosec // G115, both values fit in a byte
			pkt, _, _ := unmarshal([]byte{0x80 | byte(format), byte(packetType), 0x00, 0x00})
			typ := reflect.TypeOf(pkt)
			assert.Truef(t, corpusTypes[typ], "PT=%d FMT=%d maps to %v, not covered by the corpus", packetType, format, typ)
		}
	}
}

func FuzzParseDestinationSSRC(f *testing.F) {
	for _, datagram := range destinationSSRCCorpus(f) {
		f.Add(datagram)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		raws, rawErr := UnmarshalRaw(data)

		pkts, err := Unmarshal(data)
		if err != nil {
			// ParseDestinationSSRC must not panic or read out of bounds no
			// matter how malformed the packet.
			for _, raw := range raws {
				_, _ = raw.ParseDestinationSSRC(nil)
			}

			return
		}

		// If Unmarshal parses the packet without error, UnmarshalRaw must
		// as well, returning the same number of packets.
		if !assert.NoError(t, rawErr) || !assert.Len(t, raws, len(pkts)) {
			return
		}

		// Each packet must produce the same destination SSRCs as what
		// Unmarshal would.
		for i, pkt := range pkts {
			ssrcs, err := raws[i].ParseDestinationSSRC(nil)
			assert.NoErrorf(t, err, "packet %d (%T)", i, pkt)
			assert.ElementsMatchf(t, pkt.DestinationSSRC(), ssrcs, "packet %d (%T)", i, pkt)
		}
	})
}

func TestUnmarshalRawInvalid(t *testing.T) {
	for _, test := range []struct {
		Name string
		Data []byte
	}{
		{Name: "empty", Data: []byte{}},
		{Name: "truncated header", Data: []byte{0x80, 0xC8}},
		{Name: "bad version", Data: []byte{0x40, 0xC8, 0x00, 0x00}},
		{Name: "length overruns datagram", Data: []byte{0x80, 0xC8, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00}},
		{Name: "trailing partial packet", Data: []byte{0x80, 0xC8, 0x00, 0x00, 0x80}},
	} {
		t.Run(test.Name, func(t *testing.T) {
			_, err := UnmarshalRaw(test.Data)
			assert.Error(t, err)
		})
	}
}

// Returned packets are views into the datagram; appending to one must grow a
// copy, not overwrite the packets that follow it in the shared backing.
func TestAppendRawPacketsCapacity(t *testing.T) {
	datagram := []byte{
		0x80, 0xc3, 0x00, 0x00, // unknown packet type 195, length=0
		0x80, 0xc4, 0x00, 0x00, // unknown packet type 196, length=0
	}
	raws, err := UnmarshalRaw(datagram)
	assert.NoError(t, err)
	assert.Len(t, raws, 2)

	second := append([]byte{}, raws[1]...)
	raws[0] = append(raws[0], 0xFF)
	assert.Equal(t, second, []byte(raws[1]))
}

// The append APIs exist to be allocation-free with reused buffers; pin that
// to zero for every packet type in the corpus.
func TestRawPacketAllocations(t *testing.T) {
	var raws []RawPacket
	var ssrcs []uint32
	for i, datagram := range destinationSSRCCorpus(t) {
		// Reach steady-state buffer capacity and verify the datagram parses
		// before measuring.
		var err error
		raws, err = AppendRawPackets(raws[:0], datagram)
		assert.NoErrorf(t, err, "corpus datagram %d", i)
		for _, raw := range raws {
			ssrcs, err = raw.ParseDestinationSSRC(ssrcs[:0])
			assert.NoErrorf(t, err, "corpus datagram %d", i)
		}

		allocs := testing.AllocsPerRun(100, func() {
			raws, _ = AppendRawPackets(raws[:0], datagram)
			for _, raw := range raws {
				ssrcs, _ = raw.ParseDestinationSSRC(ssrcs[:0])
			}
		})
		assert.Zerof(t, allocs, "corpus datagram %d allocates", i)
	}
}
