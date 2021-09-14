package rtcp

import (
	"fmt"
	"reflect"
	"testing"
)

// Assert that ExtendedReport is a Packet
var _ Packet = (*ExtendedReport)(nil)

// Assert that all the extended report blocks implement the interface
var (
	_ ReportBlock = (*LossRLEReportBlock)(nil)
	_ ReportBlock = (*DuplicateRLEReportBlock)(nil)
	_ ReportBlock = (*PacketReceiptTimesReportBlock)(nil)
	_ ReportBlock = (*ReceiverReferenceTimeReportBlock)(nil)
	_ ReportBlock = (*DLRRReportBlock)(nil)
	_ ReportBlock = (*StatisticsSummaryReportBlock)(nil)
	_ ReportBlock = (*VoIPMetricsReportBlock)(nil)
	_ ReportBlock = (*UnknownReportBlock)(nil)
)

func testPacket() Packet {
	return &ExtendedReport{
		SenderSSRC: 0x01020304,
		Reports: []ReportBlock{
			&LossRLEReportBlock{
				XRHeader: XRHeader{
					BlockType: LossRLEReportBlockType,
				},
				T:        12,
				SSRC:     0x12345689,
				BeginSeq: 5,
				EndSeq:   12,
				Chunks: []Chunk{
					Chunk(0x4006),
					Chunk(0x0006),
					Chunk(0x8765),
					Chunk(0x0000),
				},
			},
			&DuplicateRLEReportBlock{
				XRHeader: XRHeader{
					BlockType: DuplicateRLEReportBlockType,
				},
				T:        6,
				SSRC:     0x12345689,
				BeginSeq: 5,
				EndSeq:   12,
				Chunks: []Chunk{
					Chunk(0x4123),
					Chunk(0x3FFF),
					Chunk(0xFFFF),
					Chunk(0x0000),
				},
			},
			&PacketReceiptTimesReportBlock{
				XRHeader: XRHeader{
					BlockType: PacketReceiptTimesReportBlockType,
				},
				T:        3,
				SSRC:     0x98765432,
				BeginSeq: 15432,
				EndSeq:   15577,
				ReceiptTime: []uint32{
					0x11111111,
					0x22222222,
					0x33333333,
					0x44444444,
					0x55555555,
				},
			},
			&ReceiverReferenceTimeReportBlock{
				XRHeader: XRHeader{
					BlockType: ReceiverReferenceTimeReportBlockType,
				},
				NTPTimestamp: 0x0102030405060708,
			},
			&DLRRReportBlock{
				XRHeader: XRHeader{
					BlockType: DLRRReportBlockType,
				},
				Reports: []DLRRReport{
					{
						SSRC:   0x88888888,
						LastRR: 0x12345678,
						DLRR:   0x99999999,
					},
					{
						SSRC:   0x09090909,
						LastRR: 0x12345678,
						DLRR:   0x99999999,
					},
					{
						SSRC:   0x11223344,
						LastRR: 0x12345678,
						DLRR:   0x99999999,
					},
				},
			},
			&StatisticsSummaryReportBlock{
				XRHeader{
					BlockType: StatisticsSummaryReportBlockType,
				},
				true, true, true, ToHIPv4,
				0xFEDCBA98,
				0x1234, 0x5678,
				0x11111111,
				0x22222222,
				0x33333333,
				0x44444444,
				0x55555555,
				0x66666666,
				0x01, 0x02, 0x03, 0x04,
			},
			&VoIPMetricsReportBlock{
				XRHeader{
					BlockType: VoIPMetricsReportBlockType,
				},
				0x89ABCDEF,
				0x05, 0x06, 0x07, 0x08,
				0x1111, 0x2222, 0x3333, 0x4444,
				0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
				0x00,
				0x1122, 0x3344, 0x5566,
			},
		},
	}
}

func encodedPacket() []byte {
	return []byte{
		// RTP Header
		0x80, 0xCF, 0x00, 0x33, // byte 0 - 3
		// Sender SSRC
		0x01, 0x02, 0x03, 0x04,
		// Loss RLE Report Block
		0x01, 0x0C, 0x00, 0x04, // byte 8 - 11
		// Source SSRC
		0x12, 0x34, 0x56, 0x89,
		// Begin & End Seq
		0x00, 0x05, 0x00, 0x0C, // byte 16 - 19
		// Chunks
		0x40, 0x06, 0x00, 0x06,
		0x87, 0x65, 0x00, 0x00, // byte 24 - 27
		// Duplicate RLE Report Block
		0x02, 0x06, 0x00, 0x04,
		// Source SSRC
		0x12, 0x34, 0x56, 0x89, // byte 32 - 35
		// Begin & End Seq
		0x00, 0x05, 0x00, 0x0C,
		// Chunks
		0x41, 0x23, 0x3F, 0xFF, // byte 40 - 43
		0xFF, 0xFF, 0x00, 0x00,
		// Packet Receipt Times Report Block
		0x03, 0x03, 0x00, 0x07, // byte 48 - 51
		// Source SSRC
		0x98, 0x76, 0x54, 0x32,
		// Begin & End Seq
		0x3C, 0x48, 0x3C, 0xD9, // byte 56 - 59
		// Receipt times
		0x11, 0x11, 0x11, 0x11,
		0x22, 0x22, 0x22, 0x22, // byte 64 - 67
		0x33, 0x33, 0x33, 0x33,
		0x44, 0x44, 0x44, 0x44, // byte 72 - 75
		0x55, 0x55, 0x55, 0x55,
		// Receiver Reference Time Report
		0x04, 0x00, 0x00, 0x02, // byte 80 - 83
		// Timestamp
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08, // byte 88 - 91
		// DLRR Report
		0x05, 0x00, 0x00, 0x09,
		// SSRC 1
		0x88, 0x88, 0x88, 0x88, // byte 96 - 99
		// LastRR 1
		0x12, 0x34, 0x56, 0x78,
		// DLRR 1
		0x99, 0x99, 0x99, 0x99, // byte 104 - 107
		// SSRC 2
		0x09, 0x09, 0x09, 0x09,
		// LastRR 2
		0x12, 0x34, 0x56, 0x78, // byte 112 - 115
		// DLRR 2
		0x99, 0x99, 0x99, 0x99,
		// SSRC 3
		0x11, 0x22, 0x33, 0x44, // byte 120 - 123
		// LastRR 3
		0x12, 0x34, 0x56, 0x78,
		// DLRR 3
		0x99, 0x99, 0x99, 0x99, // byte 128 - 131
		// Statistics Summary Report
		0x06, 0xE8, 0x00, 0x09,
		// SSRC
		0xFE, 0xDC, 0xBA, 0x98, // byte 136 - 139
		// Various statistics
		0x12, 0x34, 0x56, 0x78,
		0x11, 0x11, 0x11, 0x11, // byte 144 - 147
		0x22, 0x22, 0x22, 0x22,
		0x33, 0x33, 0x33, 0x33, // byte 152 - 155
		0x44, 0x44, 0x44, 0x44,
		0x55, 0x55, 0x55, 0x55, // byte 160 - 163
		0x66, 0x66, 0x66, 0x66,
		0x01, 0x02, 0x03, 0x04, // byte 168 - 171
		// VoIP Metrics Report
		0x07, 0x00, 0x00, 0x08,
		// SSRC
		0x89, 0xAB, 0xCD, 0xEF, // byte 176 - 179
		// Various statistics
		0x05, 0x06, 0x07, 0x08,
		0x11, 0x11, 0x22, 0x22, // byte 184 - 187
		0x33, 0x33, 0x44, 0x44,
		0x11, 0x22, 0x33, 0x44, // byte 192 - 195
		0x55, 0x66, 0x77, 0x88,
		0x99, 0x00, 0x11, 0x22, // byte 200 - 203
		0x33, 0x44, 0x55, 0x66, // byte 204 - 207
	}
}

func TestEncode(t *testing.T) {
	expected := encodedPacket()
	packet := testPacket()
	rawPacket, err := packet.Marshal()
	if err != nil {
		t.Fatalf("Error marshaling packet: %v", err)
	}
	if len(rawPacket) != len(expected) {
		t.Fatalf("Encoded message is %d bytes; expected is %d", len(rawPacket), len(expected))
	}

	for i := 0; i < len(rawPacket); i++ {
		if rawPacket[i] != expected[i] {
			t.Errorf("Byte %d of encoded packet does not match: expected 0x%02X, got 0x%02X", i, expected[i], rawPacket[i])
		}
	}
}

func TestDecode(t *testing.T) {
	encoded := encodedPacket()
	expected := testPacket()

	// We need to make sure the header has been set up correctly
	// before we test for equality
	for _, p := range expected.(*ExtendedReport).Reports {
		p.setupBlockHeader()
	}

	p := new(ExtendedReport)
	err := p.Unmarshal(encoded)
	if err != nil {
		t.Fatalf("Error unmarshaling packet: %v", err)
	}

	if !reflect.DeepEqual(p, expected) {
		t.Errorf("(deep equal) Decoded packet does not match expected packet")
	}

	if p.String() != expected.(fmt.Stringer).String() {
		t.Errorf("(string compare) Decoded packet does not match expected packet")
	}
}
