package rtcp

import (
	"testing"
)

func TestPrint(t *testing.T) {
	type Tests struct {
		packet   Packet
		expected string
	}

	tests := []Tests{
		{
			&ExtendedReport{
				SenderSSRC: 0x01020304,
				Reports: []ReportBlock{
					&LossRLEReportBlock{
						XRHeader: XRHeader{
							BlockType: LossRLEReportBlockType,
						},
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
			},
			"rtcp.ExtendedReport:\n" +
				"\tSenderSSRC: 0x1020304\n" +
				"\tReports:\n" +
				"\t\t0 (rtcp.LossRLEReportBlock):\n" +
				"\t\t\tXRHeader:\n" +
				"\t\t\t\tBlockType: [LossRLEReportBlockType]\n" +
				"\t\t\t\tTypeSpecific: 0x0\n" +
				"\t\t\t\tBlockLength: 0\n" +
				"\t\t\tT: 0\n" +
				"\t\t\tSSRC: 0x12345689\n" +
				"\t\t\tBeginSeq: 5\n" +
				"\t\t\tEndSeq: 12\n" +
				"\t\t\tChunks:\n" +
				"\t\t\t\t0: [[RunLength type=1, length=6]]\n" +
				"\t\t\t\t1: [[RunLength type=0, length=6]]\n" +
				"\t\t\t\t2: [[BitVector 0b000011101100101]]\n" +
				"\t\t\t\t3: [[TerminatingNull]]\n" +
				"\t\t1 (rtcp.DuplicateRLEReportBlock):\n" +
				"\t\t\tXRHeader:\n" +
				"\t\t\t\tBlockType: [DuplicateRLEReportBlockType]\n" +
				"\t\t\t\tTypeSpecific: 0x0\n" +
				"\t\t\t\tBlockLength: 0\n" +
				"\t\t\tT: 0\n" +
				"\t\t\tSSRC: 0x12345689\n" +
				"\t\t\tBeginSeq: 5\n" +
				"\t\t\tEndSeq: 12\n" +
				"\t\t\tChunks:\n" +
				"\t\t\t\t0: [[RunLength type=1, length=291]]\n" +
				"\t\t\t\t1: [[RunLength type=0, length=16383]]\n" +
				"\t\t\t\t2: [[BitVector 0b111111111111111]]\n" +
				"\t\t\t\t3: [[TerminatingNull]]\n" +
				"\t\t2 (rtcp.PacketReceiptTimesReportBlock):\n" +
				"\t\t\tXRHeader:\n" +
				"\t\t\t\tBlockType: [PacketReceiptTimesReportBlockType]\n" +
				"\t\t\t\tTypeSpecific: 0x0\n" +
				"\t\t\t\tBlockLength: 0\n" +
				"\t\t\tT: 0\n" +
				"\t\t\tSSRC: 0x98765432\n" +
				"\t\t\tBeginSeq: 15432\n" +
				"\t\t\tEndSeq: 15577\n" +
				"\t\t\tReceiptTime: [286331153 572662306 858993459 1145324612 1431655765]\n" +
				"\t\t3 (rtcp.ReceiverReferenceTimeReportBlock):\n" +
				"\t\t\tXRHeader:\n" +
				"\t\t\t\tBlockType: [ReceiverReferenceTimeReportBlockType]\n" +
				"\t\t\t\tTypeSpecific: 0x0\n" +
				"\t\t\t\tBlockLength: 0\n" +
				"\t\t\tNTPTimestamp: 72623859790382856\n" +
				"\t\t4 (rtcp.DLRRReportBlock):\n" +
				"\t\t\tXRHeader:\n" +
				"\t\t\t\tBlockType: [DLRRReportBlockType]\n" +
				"\t\t\t\tTypeSpecific: 0x0\n" +
				"\t\t\t\tBlockLength: 0\n" +
				"\t\t\tReports:\n" +
				"\t\t\t\t0:\n" +
				"\t\t\t\t\tSSRC: 0x88888888\n" +
				"\t\t\t\t\tLastRR: 305419896\n" +
				"\t\t\t\t\tDLRR: 2576980377\n" +
				"\t\t\t\t1:\n" +
				"\t\t\t\t\tSSRC: 0x9090909\n" +
				"\t\t\t\t\tLastRR: 305419896\n" +
				"\t\t\t\t\tDLRR: 2576980377\n" +
				"\t\t\t\t2:\n" +
				"\t\t\t\t\tSSRC: 0x11223344\n" +
				"\t\t\t\t\tLastRR: 305419896\n" +
				"\t\t\t\t\tDLRR: 2576980377\n" +
				"\t\t5 (rtcp.StatisticsSummaryReportBlock):\n" +
				"\t\t\tXRHeader:\n" +
				"\t\t\t\tBlockType: [StatisticsSummaryReportBlockType]\n" +
				"\t\t\t\tTypeSpecific: 0x0\n" +
				"\t\t\t\tBlockLength: 0\n" +
				"\t\t\tLossReports: true\n" +
				"\t\t\tDuplicateReports: true\n" +
				"\t\t\tJitterReports: true\n" +
				"\t\t\tTTLorHopLimit: [[ToH = IPv4]]\n" +
				"\t\t\tSSRC: 0xFEDCBA98\n" +
				"\t\t\tBeginSeq: 4660\n" +
				"\t\t\tEndSeq: 22136\n" +
				"\t\t\tLostPackets: 286331153\n" +
				"\t\t\tDupPackets: 572662306\n" +
				"\t\t\tMinJitter: 858993459\n" +
				"\t\t\tMaxJitter: 1145324612\n" +
				"\t\t\tMeanJitter: 1431655765\n" +
				"\t\t\tDevJitter: 1717986918\n" +
				"\t\t\tMinTTLOrHL: 1\n" +
				"\t\t\tMaxTTLOrHL: 2\n" +
				"\t\t\tMeanTTLOrHL: 3\n" +
				"\t\t\tDevTTLOrHL: 4\n" +
				"\t\t6 (rtcp.VoIPMetricsReportBlock):\n" +
				"\t\t\tXRHeader:\n" +
				"\t\t\t\tBlockType: [VoIPMetricsReportBlockType]\n" +
				"\t\t\t\tTypeSpecific: 0x0\n" +
				"\t\t\t\tBlockLength: 0\n" +
				"\t\t\tSSRC: 0x89ABCDEF\n" +
				"\t\t\tLossRate: 5\n" +
				"\t\t\tDiscardRate: 6\n" +
				"\t\t\tBurstDensity: 7\n" +
				"\t\t\tGapDensity: 8\n" +
				"\t\t\tBurstDuration: 4369\n" +
				"\t\t\tGapDuration: 8738\n" +
				"\t\t\tRoundTripDelay: 13107\n" +
				"\t\t\tEndSystemDelay: 17476\n" +
				"\t\t\tSignalLevel: 17\n" +
				"\t\t\tNoiseLevel: 34\n" +
				"\t\t\tRERL: 51\n" +
				"\t\t\tGmin: 68\n" +
				"\t\t\tRFactor: 85\n" +
				"\t\t\tExtRFactor: 102\n" +
				"\t\t\tMOSLQ: 119\n" +
				"\t\t\tMOSCQ: 136\n" +
				"\t\t\tRXConfig: 153\n" +
				"\t\t\tJBNominal: 4386\n" +
				"\t\t\tJBMaximum: 13124\n" +
				"\t\t\tJBAbsMax: 21862\n",
		},
		{
			&FullIntraRequest{
				SenderSSRC: 0x0,
				MediaSSRC:  0x4bc4fcb4,
				FIR: []FIREntry{
					{
						SSRC:           0x12345678,
						SequenceNumber: 0x42,
					},
					{
						SSRC:           0x98765432,
						SequenceNumber: 0x57,
					},
				},
			},
			"rtcp.FullIntraRequest:\n" +
				"\tSenderSSRC: 0\n" +
				"\tMediaSSRC: 1271200948\n" +
				"\tFIR:\n" +
				"\t\t0:\n" +
				"\t\t\tSSRC: 305419896\n" +
				"\t\t\tSequenceNumber: 66\n" +
				"\t\t1:\n" +
				"\t\t\tSSRC: 2557891634\n" +
				"\t\t\tSequenceNumber: 87\n",
		},
		{
			&Goodbye{
				Sources: []uint32{
					0x01020304,
					0x05060708,
				},
				Reason: "because",
			},
			"rtcp.Goodbye:\n" +
				"\tSources: [16909060 84281096]\n" +
				"\tReason: because\n",
		},
		{
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
			"rtcp.ReceiverReport:\n" +
				"\tSSRC: 2419039790\n" +
				"\tReports:\n" +
				"\t\t0:\n" +
				"\t\t\tSSRC: 3160316480\n" +
				"\t\t\tFractionLost: 0\n" +
				"\t\t\tTotalLost: 0\n" +
				"\t\t\tLastSequenceNumber: 18145\n" +
				"\t\t\tJitter: 273\n" +
				"\t\t\tLastSenderReport: 166945842\n" +
				"\t\t\tDelay: 150137\n" +
				"\tProfileExtensions: []\n",
		},
		{
			NewCNAMESourceDescription(0x902f9e2e, "{9c00eb92-1afb-9d49-a47d-91f64eee69f5}"),
			"rtcp.SourceDescription:\n" +
				"\tChunks:\n" +
				"\t\t0:\n" +
				"\t\t\tSource: 2419039790\n" +
				"\t\t\tItems:\n" +
				"\t\t\t\t0:\n" +
				"\t\t\t\t\tType: [CNAME]\n" +
				"\t\t\t\t\tText: {9c00eb92-1afb-9d49-a47d-91f64eee69f5}\n",
		},
		{
			&PictureLossIndication{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
			},
			"rtcp.PictureLossIndication:\n" +
				"\tSenderSSRC: 2419039790\n" +
				"\tMediaSSRC: 2419039790\n",
		},
		{
			&RapidResynchronizationRequest{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
			},
			"rtcp.RapidResynchronizationRequest:\n" +
				"\tSenderSSRC: 2419039790\n" +
				"\tMediaSSRC: 2419039790\n",
		},
		{
			&ReceiverEstimatedMaximumBitrate{
				SenderSSRC: 1,
				Bitrate:    8927168,
				SSRCs:      []uint32{1215622422},
			},
			"rtcp.ReceiverEstimatedMaximumBitrate:\n" +
				"\tSenderSSRC: 1\n" +
				"\tBitrate: 8.927168e+06\n" +
				"\tSSRCs: [1215622422]\n",
		},
		{
			&SenderReport{
				SSRC:        0x902f9e2e,
				NTPTime:     0xda8bd1fcdddda05a,
				RTPTime:     0xaaf4edd5,
				PacketCount: 1,
				OctetCount:  2,
				Reports: []ReceptionReport{{
					SSRC:               0xbc5e9a40,
					FractionLost:       0,
					TotalLost:          0,
					LastSequenceNumber: 0x46e1,
					Jitter:             273,
					LastSenderReport:   0x9f36432,
					Delay:              150137,
				}},
				ProfileExtensions: []byte{
					0x81, 0xca, 0x0, 0x6,
					0x2b, 0x7e, 0xc0, 0xc5,
					0x1, 0x10, 0x4c, 0x63,
					0x49, 0x66, 0x7a, 0x58,
					0x6f, 0x6e, 0x44, 0x6f,
					0x72, 0x64, 0x53, 0x65,
					0x57, 0x36, 0x0, 0x0,
				},
			},
			"rtcp.SenderReport:\n" +
				"\tSSRC: 2419039790\n" +
				"\tNTPTime: 15747911406015324250\n" +
				"\tRTPTime: 2868178389\n" +
				"\tPacketCount: 1\n" +
				"\tOctetCount: 2\n" +
				"\tReports:\n" +
				"\t\t0:\n" +
				"\t\t\tSSRC: 3160316480\n" +
				"\t\t\tFractionLost: 0\n" +
				"\t\t\tTotalLost: 0\n" +
				"\t\t\tLastSequenceNumber: 18145\n" +
				"\t\t\tJitter: 273\n" +
				"\t\t\tLastSenderReport: 166945842\n" +
				"\t\t\tDelay: 150137\n" +
				"\tProfileExtensions: [129 202 0 6 43 126 192 197 1 16 76 99 73 102 122 88 111 110 68 111 114 100 83 101 87 54 0 0]\n",
		},
		{
			&SliceLossIndication{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
				SLI:        []SLIEntry{{0xaaa, 0, 0x2C}},
			},
			"rtcp.SliceLossIndication:\n" +
				"\tSenderSSRC: 2419039790\n" +
				"\tMediaSSRC: 2419039790\n" +
				"\tSLI:\n" +
				"\t\t0:\n" +
				"\t\t\tFirst: 2730\n" +
				"\t\t\tNumber: 0\n" +
				"\t\t\tPicture: 44\n",
		},
		{
			&SourceDescription{
				Chunks: []SourceDescriptionChunk{
					{
						Source: 0x10000000,
						Items: []SourceDescriptionItem{
							{
								Type: SDESCNAME,
								Text: "A",
							},
							{
								Type: SDESPhone,
								Text: "B",
							},
						},
					},
				},
			},
			"rtcp.SourceDescription:\n" +
				"\tChunks:\n" +
				"\t\t0:\n" +
				"\t\t\tSource: 268435456\n" +
				"\t\t\tItems:\n" +
				"\t\t\t\t0:\n" +
				"\t\t\t\t\tType: [CNAME]\n" +
				"\t\t\t\t\tText: A\n" +
				"\t\t\t\t1:\n" +
				"\t\t\t\t\tType: [PHONE]\n" +
				"\t\t\t\t\tText: B\n",
		},
		{
			&TransportLayerCC{
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
			"rtcp.TransportLayerCC:\n" +
				"\tHeader:\n" +
				"\t\tPadding: true\n" +
				"\t\tCount: 15\n" +
				"\t\tType: [TSFB]\n" +
				"\t\tLength: 5\n" +
				"\tSenderSSRC: 4195875351\n" +
				"\tMediaSSRC: 1124282272\n" +
				"\tBaseSequenceNumber: 153\n" +
				"\tPacketStatusCount: 1\n" +
				"\tReferenceTime: 4057090\n" +
				"\tFbPktCount: 23\n" +
				"\tPacketChunks:\n" +
				"\t\t0 (rtcp.RunLengthChunk):\n" +
				"\t\t\tPacketStatusChunk: <nil>\n" +
				"\t\t\tType: 0\n" +
				"\t\t\tPacketStatusSymbol: 1\n" +
				"\t\t\tRunLength: 1\n" +
				"\tRecvDeltas:\n" +
				"\t\t0:\n" +
				"\t\t\tType: 1\n" +
				"\t\t\tDelta: 37000\n",
		},
		{
			&TransportLayerNack{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
				Nacks:      []NackPair{{1, 0xAA}, {1034, 0x05}},
			},
			"rtcp.TransportLayerNack:\n" +
				"\tSenderSSRC: 2419039790\n" +
				"\tMediaSSRC: 2419039790\n" +
				"\tNacks:\n" +
				"\t\t0:\n" +
				"\t\t\tPacketID: 1\n" +
				"\t\t\tLostPackets: 170\n" +
				"\t\t1:\n" +
				"\t\t\tPacketID: 1034\n" +
				"\t\t\tLostPackets: 5\n",
		},
	}

	for i, test := range tests {
		actual := stringify(test.packet)
		if actual != test.expected {
			t.Fatalf("Error stringifying test %d\nExpected:\n%s\n\nGot:\n%s\n\n", i, test.expected, actual)
		}
	}
}
