package rtcp

import (
	"bytes"
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
