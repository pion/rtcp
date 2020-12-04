package rtcp

import (
	"errors"
	"reflect"
	"testing"
)

func TestTransportLayerNackUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      TransportLayerNack
		WantError error
	}{
		{
			Name: "valid",
			Data: []byte{
				// TransportLayerNack
				0x81, 0xcd, 0x0, 0x3,
				// sender=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// media=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// nack 0xAAAA, 0x5555
				0xaa, 0xaa, 0x55, 0x55,
			},
			Want: TransportLayerNack{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
				Nacks:      []NackPair{{0xaaaa, 0x5555}},
			},
		},
		{
			Name: "short report",
			Data: []byte{
				0x81, 0xcd, 0x0, 0x2,
				// ssrc=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// report ends early
			},
			WantError: errPacketTooShort,
		},
		{
			Name: "wrong type",
			Data: []byte{
				// v=2, p=0, count=1, SR, len=7
				0x81, 0xc8, 0x0, 0x7,
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
			},
			WantError: errWrongType,
		},
		{
			Name:      "nil",
			Data:      nil,
			WantError: errPacketTooShort,
		},
	} {
		var tln TransportLayerNack
		err := tln.Unmarshal(test.Data)
		if got, want := err, test.WantError; !errors.Is(got, want) {
			t.Fatalf("Unmarshal %q rr: err = %v, want %v", test.Name, got, want)
		}
		if err != nil {
			continue
		}

		if got, want := tln, test.Want; !reflect.DeepEqual(got, want) {
			t.Fatalf("Unmarshal %q rr: got %v, want %v", test.Name, got, want)
		}
	}
}

func TestTransportLayerNackRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Report    TransportLayerNack
		WantError error
	}{
		{
			Name: "valid",
			Report: TransportLayerNack{
				SenderSSRC: 0x902f9e2e,
				MediaSSRC:  0x902f9e2e,
				Nacks:      []NackPair{{1, 0xAA}, {1034, 0x05}},
			},
		},
	} {
		data, err := test.Report.Marshal()
		if got, want := err, test.WantError; !errors.Is(got, want) {
			t.Fatalf("Marshal %q: err = %v, want %v", test.Name, got, want)
		}
		if err != nil {
			continue
		}

		var decoded TransportLayerNack
		if err := decoded.Unmarshal(data); err != nil {
			t.Fatalf("Unmarshal %q: %v", test.Name, err)
		}

		if got, want := decoded, test.Report; !reflect.DeepEqual(got, want) {
			t.Fatalf("%q tln round trip: got %#v, want %#v", test.Name, got, want)
		}
	}
}

func testNackPair(t *testing.T, s []uint16, n NackPair) {
	l := n.PacketList()
	if !reflect.DeepEqual(l, s) {
		t.Errorf("%v: expected %v, got %v", n, s, l)
	}
}

func TestNackPair(t *testing.T) {
	testNackPair(t, []uint16{42}, NackPair{42, 0})
	testNackPair(t, []uint16{42, 43}, NackPair{42, 1})
	testNackPair(t, []uint16{42, 44}, NackPair{42, 2})
	testNackPair(t, []uint16{42, 43, 44}, NackPair{42, 3})
	testNackPair(t, []uint16{42, 42 + 16}, NackPair{42, 0x8000})
}

func TestNackPairRange(t *testing.T) {
	n := NackPair{42, 2}

	out := make([]uint16, 0)
	n.Range(func(s uint16) bool {
		out = append(out, s)
		return true
	})
	if !reflect.DeepEqual(out, []uint16{42, 44}) {
		t.Errorf("Got %v", out)
	}

	out = make([]uint16, 0)
	n.Range(func(s uint16) bool {
		out = append(out, s)
		return false
	})
	if !reflect.DeepEqual(out, []uint16{42}) {
		t.Errorf("Got %v", out)
	}
}

func TestTransportLayerNackPairGeneration(t *testing.T) {
	for _, test := range []struct {
		Name            string
		SequenceNumbers []uint16
		Expected        []NackPair
	}{
		{
			"No Sequence Numbers",
			[]uint16{},
			[]NackPair{},
		},
		{
			"Single Sequence Number",
			[]uint16{100},
			[]NackPair{
				{PacketID: 100, LostPackets: 0x0},
			},
		},
		{
			"Multiple in range, Single NACKPair",
			[]uint16{100, 101, 105, 115},
			[]NackPair{
				{PacketID: 100, LostPackets: 0x4011},
			},
		},
		{
			"Multiple Ranges, Multiple NACKPair",
			[]uint16{100, 117, 500, 501, 502},
			[]NackPair{
				{PacketID: 100, LostPackets: 0},
				{PacketID: 117, LostPackets: 0},
				{PacketID: 500, LostPackets: 0x3},
			},
		},
	} {
		actual := NackPairsFromSequenceNumbers(test.SequenceNumbers)
		if !reflect.DeepEqual(actual, test.Expected) {
			t.Fatalf("%q NackPair generation mismatch: got %#v, want %#v", test.Name, actual, test.Expected)
		}
	}
}
