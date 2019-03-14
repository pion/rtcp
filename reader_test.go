package rtcp

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// An RTCP packet from a packet dump
var realPacket = []byte{
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
}

func TestUnmarshal(t *testing.T) {
	decoder := NewDecoder(bytes.NewReader(realPacket))
	packets, err := Unmarshal(realPacket)

	if err != nil {
		t.Fatalf("Error unmarshalling packets: %s", err)
	}

	if len(packets) != 5 {
		t.Fatalf("Read the wrong number of packets from input array")
	}

	// ReceiverReport
	packet := packets[0]
	decoded, err := decoder.DecodePacket()
	assert.NoError(t, err)
	assert.IsType(t, packet, (*ReceiverReport)(nil), "Unmarshalled to incorrect type")

	wantRR := &ReceiverReport{
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
	}
	if want, got := wantRR, packet; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal rr: got %#v, want %#v", got, want)
	}
	if want, got := wantRR, decoded; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal rr: got %#v, want %#v", got, want)
	}

	// SourceDescription
	packet = packets[1]
	decoded, err = decoder.DecodePacket()
	assert.NoError(t, err)
	assert.IsType(t, packet, (*SourceDescription)(nil), "Unmarshalled to incorrect type")

	wantSdes := &SourceDescription{
		Chunks: []SourceDescriptionChunk{
			{
				Source: 0x902f9e2e,
				Items: []SourceDescriptionItem{
					{
						Type: SDESCNAME,
						Text: "{9c00eb92-1afb-9d49-a47d-91f64eee69f5}",
					},
				},
			},
		},
	}

	if got, want := packet, wantSdes; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal sdes: got %#v, want %#v", got, want)
	}
	if got, want := decoded, wantSdes; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal sdes: got %#v, want %#v", got, want)
	}

	// Goodbye
	packet = packets[2]
	decoded, err = decoder.DecodePacket()
	assert.NoError(t, err)
	assert.IsType(t, packet, (*Goodbye)(nil), "Unmarshalled to incorrect type")

	wantBye := &Goodbye{
		Sources: []uint32{0x902f9e2e},
	}
	if got, want := packet, wantBye; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal bye: got %#v, want %#v", got, want)
	}
	if got, want := decoded, wantBye; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal bye: got %#v, want %#v", got, want)
	}

	// PictureLossIndication
	packet = packets[3]
	decoded, err = decoder.DecodePacket()
	assert.NoError(t, err)
	assert.IsType(t, packet, (*PictureLossIndication)(nil), "Unmarshalled to incorrect type")

	wantPli := &PictureLossIndication{
		SenderSSRC: 0x902f9e2e,
		MediaSSRC:  0x902f9e2e,
	}
	if got, want := packet, wantPli; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal pli: got %#v, want %#v", got, want)
	}
	if got, want := decoded, wantPli; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal pli: got %#v, want %#v", got, want)
	}

	// RapidResynchronizationRequest
	packet = packets[4]
	decoded, err = decoder.DecodePacket()
	assert.NoError(t, err)
	assert.IsType(t, packet, (*RapidResynchronizationRequest)(nil), "Unmarshalled to incorrect type")

	wantRrr := &RapidResynchronizationRequest{
		SenderSSRC: 0x902f9e2e,
		MediaSSRC:  0x902f9e2e,
	}
	if got, want := packet, wantRrr; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal rrr: got %#v, want %#v", got, want)
	}
	if got, want := decoded, wantRrr; !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal rrr: got %#v, want %#v", got, want)
	}
}

func TestReadEOF(t *testing.T) {
	shortHeader := []byte{
		0x81, 0xc9, // missing type & len
	}

	_, err := Unmarshal(shortHeader)
	assert.Error(t, err)

	d := NewDecoder(bytes.NewReader(shortHeader))
	_, err = d.DecodePacket()

	assert.Error(t, err)
}

func TestBadCompound(t *testing.T) {
	//trailing data!
	badcompound := realPacket[:34]
	packets, err := Unmarshal(badcompound)
	assert.Error(t, err)

	assert.Nil(t, packets)

	//illegal start -- this should return an error, but also 2 parsed packets
	//it violates the "must start with RR or SR" rule
	badcompound = realPacket[84:104]
	packets, err = Unmarshal(badcompound)
	assert.Error(t, err)
	assert.Equal(t, len(packets), 2)
	assert.Equal(t, packets[0].Header().Type, TypeGoodbye)
	assert.Equal(t, packets[1].Header().Type, TypePayloadSpecificFeedback)
}
