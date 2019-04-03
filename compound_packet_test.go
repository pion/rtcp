package rtcp

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestUnmarshalNil(t *testing.T) {
	_, err := Unmarshal(nil)
	if got, want := err, errInvalidHeader; got != want {
		t.Fatalf("Unmarshal(nil) err = %v, want %v", got, want)
	}
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

func TestValidPacket(t *testing.T) {
	cname := &SourceDescription{
		Chunks: []SourceDescriptionChunk{{
			Source: 1234,
			Items: []SourceDescriptionItem{{
				Type: SDESCNAME,
				Text: "cname",
			}},
		}},
	}

	for _, test := range []struct {
		Name   string
		Packet CompoundPacket
		Valid  bool
	}{
		{
			Name:   "empty",
			Packet: CompoundPacket{},
			Valid:  false,
		},
		{
			Name: "no cname",
			Packet: CompoundPacket{
				&SenderReport{},
			},
			Valid: false,
		},
		{
			Name: "SDES / no cname",
			Packet: CompoundPacket{
				&SenderReport{},
				&SourceDescription{},
			},
			Valid: false,
		},
		{
			Name: "just SR",
			Packet: CompoundPacket{
				&SenderReport{},
				cname,
			},
			Valid: true,
		},
		{
			Name: "multiple SRs",
			Packet: CompoundPacket{
				&SenderReport{},
				&SenderReport{},
				cname,
			},
			Valid: false,
		},
		{
			Name: "just RR",
			Packet: CompoundPacket{
				&ReceiverReport{},
				cname,
			},
			Valid: true,
		},
		{
			Name: "multiple RRs",
			Packet: CompoundPacket{
				&ReceiverReport{},
				&ReceiverReport{},
				cname,
			},
			Valid: true,
		},
		{
			Name: "goodbye",
			Packet: CompoundPacket{
				&ReceiverReport{},
				cname,
				&Goodbye{},
			},
			Valid: true,
		},
	} {
		if got, want := test.Packet.Valid(), test.Valid; got != want {
			t.Fatalf("Valid(%s) = %v, want %v", test.Name, got, want)
		}
	}
}

func TestCompoundPacketRoundTrip(t *testing.T) {
	cname := &SourceDescription{
		Chunks: []SourceDescriptionChunk{{
			Source: 1234,
			Items: []SourceDescriptionItem{{
				Type: SDESCNAME,
				Text: "cname",
			}},
		}},
	}

	for _, test := range []struct {
		Name   string
		Packet CompoundPacket
		Err    error
	}{
		{
			Name: "bye",
			Packet: CompoundPacket{
				&ReceiverReport{},
				cname,
				&Goodbye{
					Sources: []uint32{1234},
				},
			},
		},
		{
			Name: "no cname",
			Packet: CompoundPacket{
				&ReceiverReport{},
			},
			Err: errInvalidCompound,
		},
	} {
		data, err := test.Packet.Marshal()
		if got, want := err, test.Err; got != want {
			t.Fatalf("Marshal(%v) err = %v, want nil", test.Name, err)
		}
		if err != nil {
			continue
		}

		result, err := Unmarshal(data)
		if err != nil {
			t.Fatalf("Unmarshal(%v) err = %v, want nil", test.Name, err)
		}

		data2, err := result.Marshal()
		if err != nil {
			t.Fatalf("Marshal(%v) err = %v, want nil", test.Name, err)
		}

		if got, want := data, data2; !reflect.DeepEqual(got, want) {
			t.Fatalf("Unmarshal(Marshal(%v)) = %v, want %v", test.Name, got, want)
		}
	}
}
