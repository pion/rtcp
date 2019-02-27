package rtcp

import (
	"reflect"
	"testing"
)

func TestSenderReportUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      SenderReport
		WantError error
	}{
		{
			Name:      "nil",
			Data:      nil,
			WantError: errPacketTooShort,
		},
		{
			Name: "valid",
			Data: []byte{
				// v=2, p=0, count=1, SR, len=7
				0x81, 0xc8, 0x0, 0x7,
				// ssrc=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// ntp=0xda8bd1fcdddda05a
				0xda, 0x8b, 0xd1, 0xfc,
				0xdd, 0xdd, 0xa0, 0x5a,
				// rtp=0xaaf4edd5
				0xaa, 0xf4, 0xed, 0xd5,
				// packetCount=1
				0x00, 0x00, 0x00, 0x01,
				// octetCount=2
				0x00, 0x00, 0x00, 0x02,
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
			},
			Want: SenderReport{
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
			},
		},
		{
			Name: "wrong type",
			Data: []byte{
				// v=2, p=0, count=1, RR, len=7
				0x81, 0xc9, 0x0, 0x7,
				// ssrc=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// ntp=0xda8bd1fcdddda05a
				0xda, 0x8b, 0xd1, 0xfc,
				0xdd, 0xdd, 0xa0, 0x5a,
				// rtp=0xaaf4edd5
				0xaa, 0xf4, 0xed, 0xd5,
				// packetCount=1
				0x00, 0x00, 0x00, 0x01,
				// octetCount=2
				0x00, 0x00, 0x00, 0x02,
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
			},
			WantError: errWrongType,
		},
		{
			Name: "bad count in header",
			Data: []byte{
				// v=2, p=0, count=1, SR, len=7
				0x82, 0xc8, 0x0, 0x7,
				// ssrc=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// ntp=0xda8bd1fcdddda05a
				0xda, 0x8b, 0xd1, 0xfc,
				0xdd, 0xdd, 0xa0, 0x5a,
				// rtp=0xaaf4edd5
				0xaa, 0xf4, 0xed, 0xd5,
				// packetCount=1
				0x00, 0x00, 0x00, 0x01,
				// octetCount=2
				0x00, 0x00, 0x00, 0x02,
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
			},
			WantError: errPacketTooShort,
		},
		{
			Name: "with extension", // issue #447
			Data: []byte{
				// v=2, p=0, count=0, SR, len=6
				0x80, 0xc8, 0x0, 0x6,
				// ssrc=0x2b7ec0c5
				0x2b, 0x7e, 0xc0, 0xc5,
				// ntp=0xe020a2a952a53fc0
				0xe0, 0x20, 0xa2, 0xa9,
				0x52, 0xa5, 0x3f, 0xc0,
				// rtp=0x2e48a552
				0x2e, 0x48, 0xa5, 0x52,
				// packetCount=70
				0x0, 0x0, 0x0, 0x46,
				// octetCount=4637
				0x0, 0x0, 0x12, 0x1d,
				// profile-specific extension
				0x81, 0xca, 0x0, 0x6,
				0x2b, 0x7e, 0xc0, 0xc5,
				0x1, 0x10, 0x4c, 0x63,
				0x49, 0x66, 0x7a, 0x58,
				0x6f, 0x6e, 0x44, 0x6f,
				0x72, 0x64, 0x53, 0x65,
				0x57, 0x36, 0x0, 0x0,
			},
			Want: SenderReport{
				SSRC:        0x2b7ec0c5,
				NTPTime:     0xe020a2a952a53fc0,
				RTPTime:     0x2e48a552,
				PacketCount: 70,
				OctetCount:  4637,
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
			WantError: nil,
		},
	} {
		var sr SenderReport
		err := sr.Unmarshal(test.Data)
		if got, want := err, test.WantError; got != want {
			t.Fatalf("Unmarshal %q sr: err = %v, want %v", test.Name, got, want)
		}
		if err != nil {
			continue
		}

		if got, want := sr, test.Want; !reflect.DeepEqual(got, want) {
			t.Fatalf("Unmarshal %q sr: got %v, want %v", test.Name, got, want)
		}
	}
}

func TestSenderReportRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Report    SenderReport
		WantError error
	}{
		{
			Name: "valid",
			Report: SenderReport{
				SSRC:        1,
				NTPTime:     999,
				RTPTime:     555,
				PacketCount: 32,
				OctetCount:  11,
				Reports: []ReceptionReport{
					{
						SSRC:               2,
						FractionLost:       2,
						TotalLost:          3,
						LastSequenceNumber: 4,
						Jitter:             5,
						LastSenderReport:   6,
						Delay:              7,
					},
					{
						SSRC: 0,
					},
				},
			},
		},
		{
			Name: "also valid",
			Report: SenderReport{
				SSRC: 2,
				Reports: []ReceptionReport{
					{
						SSRC:               999,
						FractionLost:       30,
						TotalLost:          12345,
						LastSequenceNumber: 99,
						Jitter:             22,
						LastSenderReport:   92,
						Delay:              46,
					},
				},
			},
		},
		{
			Name: "extension",
			Report: SenderReport{
				SSRC: 2,
				Reports: []ReceptionReport{
					{
						SSRC:               999,
						FractionLost:       30,
						TotalLost:          12345,
						LastSequenceNumber: 99,
						Jitter:             22,
						LastSenderReport:   92,
						Delay:              46,
					},
				},
				ProfileExtensions: []byte{1, 2, 3, 4},
			},
		},
		{
			Name: "count overflow",
			Report: SenderReport{
				SSRC:    1,
				Reports: tooManyReports,
			},
			WantError: errTooManyReports,
		},
	} {
		data, err := test.Report.Marshal()
		if got, want := err, test.WantError; got != want {
			t.Fatalf("Marshal %q: err = %v, want %v", test.Name, got, want)
		}
		if err != nil {
			continue
		}

		var decoded SenderReport
		if err := decoded.Unmarshal(data); err != nil {
			t.Fatalf("Unmarshal %q: %v", test.Name, err)
		}

		if got, want := decoded, test.Report; !reflect.DeepEqual(got, want) {
			t.Fatalf("%q sr round trip: got %#v, want %#v", test.Name, got, want)
		}
	}
}
