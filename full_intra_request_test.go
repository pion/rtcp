package rtcp

import (
	"errors"
	"reflect"
	"testing"
)

func TestFullIntraRequestUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      FullIntraRequest
		WantError error
	}{
		{
			Name: "valid",
			Data: []byte{
				// v=2, p=0, FMT=4, PSFB, len=3
				0x84, 0xce, 0x00, 0x03,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x12345678
				0x12, 0x34, 0x56, 0x78,
				// Seqno=0x42
				0x42, 0x00, 0x00, 0x00,
			},
			Want: FullIntraRequest{
				SenderSSRC: 0x0,
				MediaSSRC:  0x4bc4fcb4,
				FIR: []FIREntry{
					{
						SSRC:           0x12345678,
						SequenceNumber: 0x42,
					},
				},
			},
		},
		{
			Name: "also valid",
			Data: []byte{
				// v=2, p=0, FMT=4, PSFB, len=3
				0x84, 0xce, 0x00, 0x05,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x12345678
				0x12, 0x34, 0x56, 0x78,
				// Seqno=0x42
				0x42, 0x00, 0x00, 0x00,
				// ssrc=0x98765432
				0x98, 0x76, 0x54, 0x32,
				// Seqno=0x57
				0x57, 0x00, 0x00, 0x00,
			},
			Want: FullIntraRequest{
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
		},
		{
			Name: "packet too short",
			Data: []byte{
				0x00, 0x00, 0x00, 0x00,
			},
			WantError: errPacketTooShort,
		},
		{
			Name: "invalid header",
			Data: []byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			WantError: errBadVersion,
		},
		{
			Name: "wrong type",
			Data: []byte{
				// v=2, p=0, FMT=4, RR, len=3
				0x84, 0xc9, 0x00, 0x03,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x12345678
				0x12, 0x34, 0x56, 0x78,
				// Seqno=0x42
				0x42, 0x00, 0x00, 0x00,
			},
			WantError: errWrongType,
		},
		{
			Name: "wrong fmt",
			Data: []byte{
				// v=2, p=0, FMT=2, PSFB, len=3
				0x82, 0xce, 0x00, 0x03,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				// ssrc=0x12345678
				0x12, 0x34, 0x56, 0x78,
				// Seqno=0x42
				0x42, 0x00, 0x00, 0x00,
			},
			WantError: errWrongType,
		},
	} {
		var fir FullIntraRequest
		err := fir.Unmarshal(test.Data)
		if got, want := err, test.WantError; !errors.Is(got, want) {
			t.Fatalf("Unmarshal %q rr: err = %v, want %v", test.Name, got, want)
		}
		if err != nil {
			continue
		}

		if got, want := fir, test.Want; !reflect.DeepEqual(got, want) {
			t.Fatalf("Unmarshal %q rr: got %v, want %v", test.Name, got, want)
		}
	}
}

func TestFullIntraRequestRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Packet    FullIntraRequest
		WantError error
	}{
		{
			Name: "valid",
			Packet: FullIntraRequest{
				SenderSSRC: 1,
				MediaSSRC:  2,
				FIR: []FIREntry{{
					SSRC:           3,
					SequenceNumber: 42,
				}},
			},
		},
		{
			Name: "also valid",
			Packet: FullIntraRequest{
				SenderSSRC: 5000,
				MediaSSRC:  6000,
				FIR: []FIREntry{{
					SSRC:           3,
					SequenceNumber: 57,
				}},
			},
		},
	} {
		data, err := test.Packet.Marshal()
		if got, want := err, test.WantError; !errors.Is(got, want) {
			t.Fatalf("Marshal %q: err = %v, want %v", test.Name, got, want)
		}
		if err != nil {
			continue
		}

		var decoded FullIntraRequest
		if err := decoded.Unmarshal(data); err != nil {
			t.Fatalf("Unmarshal %q: %v", test.Name, err)
		}

		if got, want := decoded, test.Packet; !reflect.DeepEqual(got, want) {
			t.Fatalf("%q rr round trip: got %#v, want %#v", test.Name, got, want)
		}
	}
}

func TestFullIntraRequestUnmarshalHeader(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      Header
		WantError error
	}{
		{
			Name: "valid header",
			Data: []byte{
				// v=2, p=0, FMT=1, PSFB, len=1
				0x84, 0xce, 0x00, 0x02,
				// ssrc=0x0
				0x00, 0x00, 0x00, 0x00,
				// ssrc=0x4bc4fcb4
				0x4b, 0xc4, 0xfc, 0xb4,
				0x00, 0x00, 0x00, 0x00,
			},
			Want: Header{
				Count:  FormatFIR,
				Type:   TypePayloadSpecificFeedback,
				Length: 2,
			},
		},
	} {
		var fir FullIntraRequest
		err := fir.Unmarshal(test.Data)
		if got, want := err, test.WantError; !errors.Is(got, want) {
			t.Fatalf("Unmarshal header %q rr: err = %v, want %v", test.Name, got, want)
		}
		if err != nil {
			continue
		}

		if got, want := fir.Header(), test.Want; !reflect.DeepEqual(got, want) {
			t.Fatalf("Unmarshal header %q rr: got %v, want %v", test.Name, got, want)
		}
	}
}
