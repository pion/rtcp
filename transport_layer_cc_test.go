package rtcp

import (
	"reflect"
	"testing"
)

func TestTransportLayerCC_RunLengthChunkUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      RunLengthChunk
		WantError error
	}{
		{
			// 3.1.3 example1: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example1",
			Data: []byte{0, 0xDD},
			Want: RunLengthChunk{
				Type:               TypeTCCRunLengthChunk,
				PacketStatusSymbol: TypeTCCPacketNotReceived,
				RunLength:          221,
			},
			WantError: nil,
		},
		{
			// 3.1.3 example2: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example2",
			Data: []byte{0x60, 0x18},
			Want: RunLengthChunk{
				Type:               TypeTCCRunLengthChunk,
				PacketStatusSymbol: TypeTCCPacketReceivedWithoutDelta,
				RunLength:          24,
			},
			WantError: nil,
		},
	} {
		var chunk RunLengthChunk
		err := chunk.Unmarshal(test.Data)
		if err != nil {
			t.Fatalf("Unmarshal err: %v", err)
		}
		if got, want := chunk, test.Want; got != want {
			t.Fatalf("Unmarshal %q : got = %v, want %v", test.Name, got, want)
		}
	}
}

func TestTransportLayerCC_RunLengthChunkMarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      RunLengthChunk
		Want      []byte
		WantError error
	}{
		{
			// 3.1.3 example1: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example1",
			Data: RunLengthChunk{
				Type:               TypeTCCRunLengthChunk,
				PacketStatusSymbol: TypeTCCPacketNotReceived,
				RunLength:          221,
			},
			Want:      []byte{0, 0xDD},
			WantError: nil,
		},
		{
			// 3.1.3 example2: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example2",
			Data: RunLengthChunk{
				Type:               TypeTCCRunLengthChunk,
				PacketStatusSymbol: TypeTCCPacketReceivedWithoutDelta,
				RunLength:          24,
			},
			Want:      []byte{0x60, 0x18},
			WantError: nil,
		},
	} {
		chunk := test.Data
		data, _ := chunk.Marshal()
		if got, want := data, test.Want; !reflect.DeepEqual(got, want) {
			t.Fatalf("Unmarshal %q : got = %v, want %v", test.Name, got, want)
		}
	}
}

func TestTransportLayerCC_StatusVectorChunkUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      StatusVectorChunk
		WantError error
	}{
		{
			// 3.1.4 example1: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example1",
			Data: []byte{0x9F, 0x1C},
			Want: StatusVectorChunk{
				Type:       TypeTCCStatusVectorChunk,
				SymbolSize: TypeTCCSymbolSizeOneBit,
				SymbolList: []uint16{TypeTCCPacketNotReceived, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived},
			},
			WantError: nil,
		},
		{
			// 3.1.4 example2: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example2",
			Data: []byte{0xCD, 0x50},
			Want: StatusVectorChunk{
				Type:       TypeTCCStatusVectorChunk,
				SymbolSize: TypeTCCSymbolSizeTwoBit,
				SymbolList: []uint16{TypeTCCPacketNotReceived, TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived},
			},
			WantError: nil,
		},
	} {
		var chunk StatusVectorChunk
		err := chunk.Unmarshal(test.Data)
		if err != nil {
			t.Fatalf("Unmarshal err: %v", err)
		}

		if got, want := chunk, test.Want; got.Type != want.Type || got.SymbolSize != want.SymbolSize || !reflect.DeepEqual(got.SymbolList, want.SymbolList) {
			t.Fatalf("Unmarshal %q : got = %v, want %v", test.Name, got, want)
		}
	}
}

func TestTransportLayerCC_StatusVectorChunkMarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      StatusVectorChunk
		Want      []byte
		WantError error
	}{
		{
			// 3.1.4 example1: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example1",
			Data: StatusVectorChunk{
				Type:       TypeTCCStatusVectorChunk,
				SymbolSize: TypeTCCSymbolSizeOneBit,
				SymbolList: []uint16{TypeTCCPacketNotReceived, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived},
			},
			Want:      []byte{0x9F, 0x1C},
			WantError: nil,
		},
		{
			// 3.1.4 example2: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example2",
			Data: StatusVectorChunk{
				Type:       TypeTCCStatusVectorChunk,
				SymbolSize: TypeTCCSymbolSizeTwoBit,
				SymbolList: []uint16{TypeTCCPacketNotReceived, TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived},
			},
			Want:      []byte{0xCD, 0x50},
			WantError: nil,
		},
	} {
		chunk := test.Data
		data, _ := chunk.Marshal()
		if got, want := data, test.Want; !reflect.DeepEqual(got, want) {
			t.Fatalf("Unmarshal %q : got = %v, want %v", test.Name, got, want)
		}
	}
}

func TestTransportLayerCC_RecvDeltaUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      RecvDelta
		WantError error
	}{
		{
			Name: "small delta 63.75ms",
			Data: []byte{0xFF},
			Want: RecvDelta{
				Type: TypeTCCPacketReceivedSmallDelta,
				// 255 * 250
				Delta: 63750,
			},
			WantError: nil,
		},
		{
			Name: "big delta 8191.75ms",
			Data: []byte{0x7F, 0xFF},
			Want: RecvDelta{
				Type: TypeTCCPacketReceivedLargeDelta,
				// 32767 * 250
				Delta: 8191750,
			},
			WantError: nil,
		},
		{
			Name: "big delta -8192ms",
			Data: []byte{0x80, 0x00},
			Want: RecvDelta{
				Type: TypeTCCPacketReceivedLargeDelta,
				// -32768 * 250
				Delta: -8192000,
			},
			WantError: nil,
		},
	} {
		var chunk RecvDelta
		err := chunk.Unmarshal(test.Data)
		if err != nil {
			t.Fatalf("Unmarshal err: %v", err)
		}

		if got, want := chunk, test.Want; got != want {
			t.Fatalf("Unmarshal %q : got = %v, want %v", test.Name, got, want)
		}
	}
}

func TestTransportLayerCC_RecvDeltaMarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      RecvDelta
		Want      []byte
		WantError error
	}{
		{
			Name: "small delta 63.75ms",
			Data: RecvDelta{
				Type: TypeTCCPacketReceivedSmallDelta,
				// 255 * 250
				Delta: 63750,
			},
			Want:      []byte{0xFF},
			WantError: nil,
		},
		{
			Name: "big delta 8191.75ms",
			Data: RecvDelta{
				Type: TypeTCCPacketReceivedLargeDelta,
				// 32767 * 250
				Delta: 8191750,
			},
			Want:      []byte{0x7F, 0xFF},
			WantError: nil,
		},
		{
			Name: "big delta -8192ms",
			Data: RecvDelta{
				Type: TypeTCCPacketReceivedLargeDelta,
				// -32768 * 250
				Delta: -8192000,
			},
			Want:      []byte{0x80, 0x00},
			WantError: nil,
		},
	} {
		chunk := test.Data
		data, _ := chunk.Marshal()
		if got, want := data, test.Want; !reflect.DeepEqual(got, want) {
			t.Fatalf("Unmarshal %q : got = %v, want %v", test.Name, got, want)
		}
	}
}

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |V=2|P|  FMT=15 |    PT=205     |           length              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                     SSRC of packet sender                     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      SSRC of media source                     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      base sequence number     |      packet status count      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                 reference time                | fb pkt. count |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         packet chunk          |  recv delta   |  recv delta   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// 0b10101111,0b11001101,0b00000000,0b00000101,
// 0b11111010,0b00010111,0b11111010,0b00010111,
// 0b01000011,0b00000011,0b00101111,0b10100000,
// 0b00000000,0b10011001,0b00000000,0b00000001,
// 0b00111101,0b11101000,0b00000010,0b00010111,
// 0b00100000,0b00000001,0b10010100,0b00000001,
func TestTransportLayerCC_Unmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      TransportLayerCC
		WantError error
	}{
		{
			Name: "example1",
			Data: []byte{
				0xaf, 0xcd, 0x0, 0x5,
				0xfa, 0x17, 0xfa, 0x17,
				0x43, 0x3, 0x2f, 0xa0,
				0x0, 0x99, 0x0, 0x1,
				0x3d, 0xe8, 0x2, 0x17,
				0x20, 0x1, 0x94, 0x1,
			},
			Want: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  5,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          1124282272,
				BaseSequenceNumber: 153,
				PacketStatusCount:  1,
				ReferenceTime:      4057090,
				FbPktCount:         23,
				// 0b00100000, 0b00000001
				PacketChunks: []PacketStatusChunk{
					&RunLengthChunk{
						Type:               TypeTCCRunLengthChunk,
						PacketStatusSymbol: TypeTCCPacketReceivedSmallDelta,
						RunLength:          1,
					},
				},
				// 0b10010100
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 37000,
					},
				},
			},
			WantError: nil,
		},
		{
			Name: "example2",
			Data: []byte{
				0xaf, 0xcd, 0x0, 0x6,
				0xfa, 0x17, 0xfa, 0x17,
				0x19, 0x3d, 0xd8, 0xbb,
				0x1, 0x74, 0x0, 0xe,
				0x45, 0xb1, 0x5a, 0x40,
				0xd8, 0x0, 0xf0, 0xff,
				0xd0, 0x0, 0x0, 0x3,
			},
			Want: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  6,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          423483579,
				BaseSequenceNumber: 372,
				PacketStatusCount:  14,
				ReferenceTime:      4567386,
				FbPktCount:         64,
				PacketChunks: []PacketStatusChunk{
					&StatusVectorChunk{
						Type:       TypeTCCStatusVectorChunk,
						SymbolSize: TypeTCCSymbolSizeTwoBit,
						SymbolList: []uint16{TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedLargeDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived},
					},
					&StatusVectorChunk{
						Type:       TypeTCCStatusVectorChunk,
						SymbolSize: TypeTCCSymbolSizeTwoBit,
						SymbolList: []uint16{TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketReceivedWithoutDelta},
					},
				},
				// 0b10010100
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 52000,
					},
					{
						Type:  TypeTCCPacketReceivedLargeDelta,
						Delta: 0,
					},
				},
			},
			WantError: nil,
		},
		{
			Name: "example3",
			Data: []byte{
				0xaf, 0xcd, 0x0, 0x7,
				0xfa, 0x17, 0xfa, 0x17,
				0x19, 0x3d, 0xd8, 0xbb,
				0x1, 0x74, 0x0, 0x6,
				0x45, 0xb1, 0x5a, 0x40,
				0x40, 0x2, 0x20, 0x04,
				0x1f, 0xfe, 0x1f, 0x9a,
				0xd0, 0x0, 0xd0, 0x0,
			},
			Want: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  7,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          423483579,
				BaseSequenceNumber: 372,
				PacketStatusCount:  6,
				ReferenceTime:      4567386,
				FbPktCount:         64,
				PacketChunks: []PacketStatusChunk{
					&RunLengthChunk{
						Type:               TypeTCCRunLengthChunk,
						PacketStatusSymbol: TypeTCCPacketReceivedLargeDelta,
						RunLength:          2,
					},
					&RunLengthChunk{
						Type:               TypeTCCRunLengthChunk,
						PacketStatusSymbol: TypeTCCPacketReceivedSmallDelta,
						RunLength:          4,
					},
				},
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedLargeDelta,
						Delta: 2047500,
					},
					{
						Type:  TypeTCCPacketReceivedLargeDelta,
						Delta: 2022500,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 52000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 0,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 52000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 0,
					},
				},
			},
			WantError: nil,
		},
		{
			Name: "example4",
			Data: []byte{
				0xaf, 0xcd, 0x0, 0x7,
				0xfa, 0x17, 0xfa, 0x17,
				0x19, 0x3d, 0xd8, 0xbb,
				0x0, 0x4, 0x0, 0x7,
				0x10, 0x63, 0x6e, 0x1,
				0x20, 0x7, 0x4c, 0x24,
				0x24, 0x10, 0xc, 0xc,
				0x10, 0x0, 0x0, 0x3,
			},
			Want: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  7,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          423483579,
				BaseSequenceNumber: 4,
				PacketStatusCount:  7,
				ReferenceTime:      1074030,
				FbPktCount:         1,
				PacketChunks: []PacketStatusChunk{
					&RunLengthChunk{
						Type:               TypeTCCRunLengthChunk,
						PacketStatusSymbol: TypeTCCPacketReceivedSmallDelta,
						RunLength:          7,
					},
				},
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 19000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 9000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 9000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 4000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 3000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 3000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 4000,
					},
				},
			},
			WantError: nil,
		},
		{
			Name: "example5",
			Data: []byte{
				0xaf, 0xcd, 0x0, 0x6,
				0xfa, 0x17, 0xfa, 0x17,
				0x19, 0x3d, 0xd8, 0xbb,
				0x0, 0x1, 0x0, 0xe,
				0x10, 0x63, 0x6d, 0x0,
				0xba, 0x0, 0x10, 0xc,
				0xc, 0x10, 0x0, 0x3,
			},
			Want: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  6,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          423483579,
				BaseSequenceNumber: 1,
				PacketStatusCount:  14,
				ReferenceTime:      1074029,
				FbPktCount:         0,
				PacketChunks: []PacketStatusChunk{
					&StatusVectorChunk{
						Type:       TypeTCCStatusVectorChunk,
						SymbolSize: 0,
						SymbolList: []uint16{TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived},
					},
				},
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 4000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 3000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 3000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 4000,
					},
				},
			},
			WantError: nil,
		},
		{
			Name: "example6",
			Data: []byte{
				0xaf, 0xcd, 0x0, 0x7,
				0x9b, 0x74, 0xf6, 0x1f,
				0x93, 0x71, 0xdc, 0xbc,
				0x85, 0x3c, 0x0, 0x9,
				0x63, 0xf9, 0x16, 0xb3,
				0xd5, 0x52, 0x0, 0x30,
				0x9b, 0xaa, 0x6a, 0xaa,
				0x7b, 0x1, 0x9, 0x1,
			},
			Want: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  7,
				},
				SenderSSRC:         2608133663,
				MediaSSRC:          2473712828,
				BaseSequenceNumber: 34108,
				PacketStatusCount:  9,
				ReferenceTime:      6551830,
				FbPktCount:         179,
				PacketChunks: []PacketStatusChunk{
					&StatusVectorChunk{
						Type:       TypeTCCStatusVectorChunk,
						SymbolSize: TypeTCCSymbolSizeTwoBit,
						SymbolList: []uint16{TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketReceivedLargeDelta},
					},
					&RunLengthChunk{
						Type:               TypeTCCRunLengthChunk,
						PacketStatusSymbol: TypeTCCPacketNotReceived,
						RunLength:          48,
					},
				},
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 38750,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 42500,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 26500,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 42500,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 30750,
					},
					{
						Type:  TypeTCCPacketReceivedLargeDelta,
						Delta: 66250,
					},
				},
			},
			WantError: nil,
		},
	} {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			var chunk TransportLayerCC
			err := chunk.Unmarshal(test.Data)
			if err != nil {
				t.Fatalf("Unmarshal err: %v", err)
			}
			if got, want := chunk, test.Want; !reflect.DeepEqual(got, want) {
				t.Fatalf("Unmarshal %q : got = %v, want %v", test.Name, got, want)
			}
		})
	}
}

func TestTransportLayerCC_Marshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      TransportLayerCC
		Want      []byte
		WantError error
	}{
		{
			Name: "example1",
			Data: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  5,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          1124282272,
				BaseSequenceNumber: 153,
				PacketStatusCount:  1,
				ReferenceTime:      4057090,
				FbPktCount:         23,
				// 0b00100000, 0b00000001
				PacketChunks: []PacketStatusChunk{
					&RunLengthChunk{
						Type:               TypeTCCRunLengthChunk,
						PacketStatusSymbol: TypeTCCPacketReceivedSmallDelta,
						RunLength:          1,
					},
				},
				// 0b10010100
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 37000,
					},
				},
			},
			Want: []byte{
				0xaf, 0xcd, 0x0, 0x5,
				0xfa, 0x17, 0xfa, 0x17,
				0x43, 0x3, 0x2f, 0xa0,
				0x0, 0x99, 0x0, 0x1,
				0x3d, 0xe8, 0x2, 0x17,
				0x20, 0x1, 0x94, 0x1,
			},
			WantError: nil,
		},
		{
			Name: "example2",
			Data: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  6,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          423483579,
				BaseSequenceNumber: 372,
				PacketStatusCount:  2,
				ReferenceTime:      4567386,
				FbPktCount:         64,
				PacketChunks: []PacketStatusChunk{
					&StatusVectorChunk{
						Type:       TypeTCCStatusVectorChunk,
						SymbolSize: TypeTCCSymbolSizeTwoBit,
						SymbolList: []uint16{TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedLargeDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived},
					},
					&StatusVectorChunk{
						Type:       TypeTCCStatusVectorChunk,
						SymbolSize: TypeTCCSymbolSizeTwoBit,
						SymbolList: []uint16{TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketReceivedWithoutDelta, TypeTCCPacketReceivedWithoutDelta},
					},
				},
				// 0b10010100
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 52000,
					},
					{
						Type:  TypeTCCPacketReceivedLargeDelta,
						Delta: 0,
					},
				},
			},
			Want: []byte{
				0xaf, 0xcd, 0x0, 0x6,
				0xfa, 0x17, 0xfa, 0x17,
				0x19, 0x3d, 0xd8, 0xbb,
				0x1, 0x74, 0x0, 0x2,
				0x45, 0xb1, 0x5a, 0x40,
				0xd8, 0x0, 0xf0, 0xff,
				0xd0, 0x0, 0x0, 0x2,
			},
			WantError: nil,
		},
		{
			Name: "example3",
			Data: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  7,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          423483579,
				BaseSequenceNumber: 372,
				PacketStatusCount:  6,
				ReferenceTime:      4567386,
				FbPktCount:         64,
				PacketChunks: []PacketStatusChunk{
					&RunLengthChunk{
						Type:               TypeTCCRunLengthChunk,
						PacketStatusSymbol: TypeTCCPacketReceivedLargeDelta,
						RunLength:          2,
					},
					&RunLengthChunk{
						Type:               TypeTCCRunLengthChunk,
						PacketStatusSymbol: TypeTCCPacketReceivedSmallDelta,
						RunLength:          4,
					},
				},
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedLargeDelta,
						Delta: 2047500,
					},
					{
						Type:  TypeTCCPacketReceivedLargeDelta,
						Delta: 2022500,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 52000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 0,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 52000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 0,
					},
				},
			},
			Want: []byte{
				0xaf, 0xcd, 0x0, 0x7,
				0xfa, 0x17, 0xfa, 0x17,
				0x19, 0x3d, 0xd8, 0xbb,
				0x1, 0x74, 0x0, 0x6,
				0x45, 0xb1, 0x5a, 0x40,
				0x40, 0x2, 0x20, 0x04,
				0x1f, 0xfe, 0x1f, 0x9a,
				0xd0, 0x0, 0xd0, 0x0,
			},
			WantError: nil,
		},
		{
			Name: "example4",
			Data: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  7,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          423483579,
				BaseSequenceNumber: 4,
				PacketStatusCount:  7,
				ReferenceTime:      1074030,
				FbPktCount:         1,
				PacketChunks: []PacketStatusChunk{
					&RunLengthChunk{
						Type:               TypeTCCRunLengthChunk,
						PacketStatusSymbol: TypeTCCPacketReceivedSmallDelta,
						RunLength:          7,
					},
				},
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 19000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 9000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 9000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 4000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 3000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 3000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 4000,
					},
				},
			},
			Want: []byte{
				0xaf, 0xcd, 0x0, 0x7,
				0xfa, 0x17, 0xfa, 0x17,
				0x19, 0x3d, 0xd8, 0xbb,
				0x0, 0x4, 0x0, 0x7,
				0x10, 0x63, 0x6e, 0x1,
				0x20, 0x7, 0x4c, 0x24,
				0x24, 0x10, 0xc, 0xc,
				0x10, 0x0, 0x0, 0x3,
			},
			WantError: nil,
		},
		{
			Name: "example5",
			Data: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  6,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          423483579,
				BaseSequenceNumber: 1,
				PacketStatusCount:  14,
				ReferenceTime:      1074029,
				FbPktCount:         0,
				PacketChunks: []PacketStatusChunk{
					&StatusVectorChunk{
						Type:       TypeTCCStatusVectorChunk,
						SymbolSize: TypeTCCSymbolSizeOneBit,
						SymbolList: []uint16{TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived},
					},
				},
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 4000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 3000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 3000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 4000,
					},
				},
			},
			Want: []byte{
				0xaf, 0xcd, 0x0, 0x6,
				0xfa, 0x17, 0xfa, 0x17,
				0x19, 0x3d, 0xd8, 0xbb,
				0x0, 0x1, 0x0, 0xe,
				0x10, 0x63, 0x6d, 0x0,
				0xba, 0x0, 0x10, 0xc,
				0xc, 0x10, 0x0, 0x2,
			},
			WantError: nil,
		},
		{
			Name: "example6",
			Data: TransportLayerCC{
				Header: Header{
					Padding: true,
					Count:   FormatTCC,
					Type:    TypeTransportSpecificFeedback,
					Length:  7,
				},
				SenderSSRC:         4195875351,
				MediaSSRC:          1124282272,
				BaseSequenceNumber: 39956,
				PacketStatusCount:  12,
				ReferenceTime:      7701536,
				FbPktCount:         0,
				PacketChunks: []PacketStatusChunk{
					&StatusVectorChunk{
						Type:       TypeTCCStatusVectorChunk,
						SymbolSize: TypeTCCSymbolSizeOneBit,
						SymbolList: []uint16{TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketReceivedSmallDelta, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived, TypeTCCPacketNotReceived},
					},
				},
				RecvDeltas: []*RecvDelta{
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 48250,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 15750,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 14750,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 15750,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 20750,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 36000,
					},
					{
						Type:  TypeTCCPacketReceivedSmallDelta,
						Delta: 14750,
					},
				},
			},
			Want: []byte{
				0xaf, 0xcd, 0x0, 0x7,
				0xfa, 0x17, 0xfa, 0x17,
				0x43, 0x3, 0x2f, 0xa0,
				0x9c, 0x14, 0x0, 0xc,
				0x75, 0x84, 0x20, 0x0,

				0xbe, 0xc0, 0xc1, 0x3f,
				0x3b, 0x3f, 0x53, 0x90,
				0x3b, 0x0, 0x0, 0x3,
			},
			WantError: nil,
		},
	} {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			bin, err := test.Data.Marshal()
			if err != nil {
				t.Fatalf("Marshal err: %v", err)
			}
			if got, want := bin, test.Want; !reflect.DeepEqual(got, want) {
				t.Fatalf("Marshal %q : got = %v, want %v", test.Name, got, want)
			}
		})
	}
}
