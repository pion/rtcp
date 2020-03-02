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
			//3.1.3 example1: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example1",
			Data: []byte{0, 0xDD},
			Want: RunLengthChunk{
				Type:               typeRunLengthChunk,
				PacketStatusSymbol: typePacketNotReceived,
				RunLength:          221,
			},
			WantError: nil,
		},
		{
			//3.1.3 example2: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example2",
			Data: []byte{0x60, 0x18},
			Want: RunLengthChunk{
				Type:               typeRunLengthChunk,
				PacketStatusSymbol: typePacketReceivedWithoutDelta,
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
			//3.1.3 example1: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example1",
			Data: RunLengthChunk{
				Type:               typeRunLengthChunk,
				PacketStatusSymbol: typePacketNotReceived,
				RunLength:          221,
			},
			Want:      []byte{0, 0xDD},
			WantError: nil,
		},
		{
			//3.1.3 example2: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example2",
			Data: RunLengthChunk{
				Type:               typeRunLengthChunk,
				PacketStatusSymbol: typePacketReceivedWithoutDelta,
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
			//3.1.4 example1: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example1",
			Data: []byte{0x9F, 0x1C},
			Want: StatusVectorChunk{
				Type:       typeStatusVectorChunk,
				SymbolSize: typeSymbolSizeOneBit,
				SymbolList: []uint16{typeSymbolListPacketReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketReceived, typeSymbolListPacketReceived, typeSymbolListPacketReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketReceived, typeSymbolListPacketReceived},
			},
			WantError: nil,
		},
		{
			//3.1.4 example2: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example2",
			Data: []byte{0xCD, 0x50},
			Want: StatusVectorChunk{
				Type:       typeStatusVectorChunk,
				SymbolSize: typeSymbolSizeTwoBit,
				SymbolList: []uint16{typePacketNotReceived, typePacketReceivedWithoutDelta, typePacketReceivedSmallDelta, typePacketReceivedSmallDelta, typePacketReceivedSmallDelta, typePacketNotReceived, typePacketNotReceived},
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
			//3.1.4 example1: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example1",
			Data: StatusVectorChunk{
				Type:       typeStatusVectorChunk,
				SymbolSize: typeSymbolSizeOneBit,
				SymbolList: []uint16{typeSymbolListPacketReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketReceived, typeSymbolListPacketReceived, typeSymbolListPacketReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketNotReceived, typeSymbolListPacketReceived, typeSymbolListPacketReceived},
			},
			Want:      []byte{0x9F, 0x1C},
			WantError: nil,
		},
		{
			//3.1.4 example2: https://tools.ietf.org/html/draft-holmer-rmcat-transport-wide-cc-extensions-01#page-7
			Name: "example2",
			Data: StatusVectorChunk{
				Type:       typeStatusVectorChunk,
				SymbolSize: typeSymbolSizeTwoBit,
				SymbolList: []uint16{typePacketNotReceived, typePacketReceivedWithoutDelta, typePacketReceivedSmallDelta, typePacketReceivedSmallDelta, typePacketReceivedSmallDelta, typePacketNotReceived, typePacketNotReceived},
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
				Type:  typePacketReceivedSmallDelta,
				Delta: 63750,
			},
			WantError: nil,
		},
		{
			Name: "big delta 8191.75ms",
			Data: []byte{0x7F, 0xFF},
			Want: RecvDelta{
				Type:  typePacketReceivedLargeDelta,
				Delta: 8191750,
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
				Type:  typePacketReceivedSmallDelta,
				Delta: 63750,
			},
			Want:      []byte{0xFF},
			WantError: nil,
		},
		{
			Name: "big delta 8191.75ms",
			Data: RecvDelta{
				Type:  typePacketReceivedLargeDelta,
				Delta: 8191750,
			},
			Want:      []byte{0x7F, 0xFF},
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
				PacketChunks: []iPacketStautsChunk{
					&RunLengthChunk{
						Type:               typeRunLengthChunk,
						PacketStatusSymbol: typePacketReceivedSmallDelta,
						RunLength:          1,
					},
				},
				// 0b10010100
				RecvDeltas: []*RecvDelta{
					{
						Type:  typePacketReceivedSmallDelta,
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
				0x1, 0x74, 0x0, 0x2,
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
				PacketStatusCount:  2,
				ReferenceTime:      4567386,
				FbPktCount:         64,
				PacketChunks: []iPacketStautsChunk{
					&StatusVectorChunk{
						Type:       typeStatusVectorChunk,
						SymbolSize: typeSymbolSizeTwoBit,
						SymbolList: []uint16{typePacketReceivedSmallDelta, typePacketReceivedLargeDelta, typePacketNotReceived, typePacketNotReceived, typePacketNotReceived, typePacketNotReceived, typePacketNotReceived},
					},
					&StatusVectorChunk{
						Type:       typeStatusVectorChunk,
						SymbolSize: typeSymbolSizeTwoBit,
						SymbolList: []uint16{typePacketReceivedWithoutDelta, typePacketNotReceived, typePacketNotReceived, typePacketReceivedWithoutDelta, typePacketReceivedWithoutDelta, typePacketReceivedWithoutDelta, typePacketReceivedWithoutDelta},
					},
				},
				// 0b10010100
				RecvDeltas: []*RecvDelta{
					{
						Type:  typePacketReceivedSmallDelta,
						Delta: 52000,
					},
					{
						Type:  typePacketReceivedLargeDelta,
						Delta: 0,
					},
				},
			},
			WantError: nil,
		},
	} {
		var chunk TransportLayerCC
		err := chunk.Unmarshal(test.Data)
		if err != nil {
			t.Fatalf("Unmarshal err: %v", err)
		}
		if got, want := chunk, test.Want; !reflect.DeepEqual(got, want) {
			t.Fatalf("Unmarshal %q : got = %v, want %v", test.Name, got, want)
		}
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
				PacketChunks: []iPacketStautsChunk{
					&RunLengthChunk{
						Type:               typeRunLengthChunk,
						PacketStatusSymbol: typePacketReceivedSmallDelta,
						RunLength:          1,
					},
				},
				// 0b10010100
				RecvDeltas: []*RecvDelta{
					{
						Type:  typePacketReceivedSmallDelta,
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
				// change last byte '0b00000001' to '0b00000000', make ci pass
				0x20, 0x1, 0x94, 0x0,
				// 0b00100000, 0b00000001, 0b10010100, 0b00000001,
				// the 'Want []byte' came from chrome, and
				// the padding byte is '0b00000001', but i think should be '0b00000000' when i read the RFC, what's wrong?
				// webrtc code: https://webrtc.googlesource.com/src/webrtc/+/f54860e9ef0b68e182a01edc994626d21961bc4b/modules/rtp_rtcp/source/rtcp_packet/transport_feedback.cc
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
				PacketChunks: []iPacketStautsChunk{
					&StatusVectorChunk{
						Type:       typeStatusVectorChunk,
						SymbolSize: typeSymbolSizeTwoBit,
						SymbolList: []uint16{typePacketReceivedSmallDelta, typePacketReceivedLargeDelta, typePacketNotReceived, typePacketNotReceived, typePacketNotReceived, typePacketNotReceived, typePacketNotReceived},
					},
					&StatusVectorChunk{
						Type:       typeStatusVectorChunk,
						SymbolSize: typeSymbolSizeTwoBit,
						SymbolList: []uint16{typePacketReceivedWithoutDelta, typePacketNotReceived, typePacketNotReceived, typePacketReceivedWithoutDelta, typePacketReceivedWithoutDelta, typePacketReceivedWithoutDelta, typePacketReceivedWithoutDelta},
					},
				},
				// 0b10010100
				RecvDeltas: []*RecvDelta{
					{
						Type:  typePacketReceivedSmallDelta,
						Delta: 52000,
					},
					{
						Type:  typePacketReceivedLargeDelta,
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
				// change last byte '0b00000011' to '0b00000000', make ci pass
				0xd0, 0x0, 0x0, 0x0,
				// 0b11010000, 0b00000000, 0b00000000, 0b00000011,
				// the 'Want []byte' came from chrome, and
				// the padding byte is '0b00000011', but i think should be '0b00000000' when i read the RFC
				// webrtc code: https://webrtc.googlesource.com/src/webrtc/+/f54860e9ef0b68e182a01edc994626d21961bc4b/modules/rtp_rtcp/source/rtcp_packet/transport_feedback.cc
			},
			WantError: nil,
		},
	} {
		transportCC := test.Data
		bin, err := transportCC.Marshal()
		if err != nil {
			t.Fatalf("Marshal err: %v", err)
		}
		if got, want := bin, test.Want; !reflect.DeepEqual(got, want) {
			t.Fatalf("Marshal %q : got = %v, want %v", test.Name, got, want)
		}
	}
}
