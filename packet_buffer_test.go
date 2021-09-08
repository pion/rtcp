package rtcp

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestWrite(t *testing.T) {
	type Subtree struct {
		SubA uint32
		SubB []uint8
	}

	s := struct {
		A uint8
		Z uint32 `encoding:"omit"`
		B uint16
		C uint32
		D uint64
		_ uint8
		E []uint16
		F Subtree
		G []Subtree
	}{
		0xf8,
		0x01234567,
		0x1234,
		0x56789ABC,
		0x0102030405060708,
		0x12,
		[]uint16{0x0E, 0x02FF},
		Subtree{0x11223344, []uint8{9, 8, 7, 6, 5, 4, 3, 2, 1}},
		[]Subtree{{0x01, []uint8{1, 2, 3, 4}}, {0x02, []uint8{5, 6, 7, 8}}},
	}
	expected := []byte{
		0xf8,
		0x12, 0x34,
		0x56, 0x78, 0x9A, 0xBC,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x00,
		0x00, 0x0E, 0x02, 0xFF,
		0x11, 0x22, 0x33, 0x44, 9, 8, 7, 6, 5, 4, 3, 2, 1,
		0x00, 0x00, 0x00, 0x01, 1, 2, 3, 4, 0x00, 0x00, 0x00, 0x02, 5, 6, 7, 8,
	}

	size := wireSize(s)
	if size != len(expected) {
		t.Fatalf("wireSize() returned unexpected value. Expected %v, got %v", len(expected), size)
	}

	raw := make([]byte, len(expected))
	buffer := packetBuffer{bytes: raw}
	err := buffer.write(s)
	if err != nil {
		t.Fatalf("Serialization failed. Err = %v", err)
	}
	if !bytes.Equal(raw, expected) {
		t.Fatalf("Serialization failed. Wanted %v, got %v", expected, raw)
	}

	// Check for overflow
	raw = make([]byte, len(expected)-1)
	buffer = packetBuffer{bytes: raw}
	err = buffer.write(s)
	if !errors.Is(err, errWrongMarshalSize) {
		t.Fatalf("Serialization failed. Err = %v", err)
	}
}

func TestReadUint8(t *testing.T) {
	const expected = 0x01
	raw := []byte{expected}
	output := uint8(0)
	buffer := packetBuffer{bytes: raw}
	err := buffer.read(&output)
	if err != nil {
		t.Fatalf("Value parsing failed. Err = %v", err)
	}
	if output != expected {
		t.Fatalf("Reading uint8 failed. Wanted %X, got %X", expected, output)
	}
}

func TestReadUint16(t *testing.T) {
	const expected = 0x0102
	raw := []byte{0x01, 0x02}
	output := uint16(0)
	buffer := packetBuffer{bytes: raw}
	err := buffer.read(&output)
	if err != nil {
		t.Fatalf("Value parsing failed. Err = %v", err)
	}
	if output != expected {
		t.Fatalf("Reading uint16 failed. Wanted %X, got %X", expected, output)
	}
}

func TestReadUint32(t *testing.T) {
	const expected = 0x01020304
	raw := []byte{0x01, 0x02, 0x03, 0x04}
	output := uint32(0)
	buffer := packetBuffer{bytes: raw}
	err := buffer.read(&output)
	if err != nil {
		t.Fatalf("Value parsing failed. Err = %v", err)
	}
	if output != expected {
		t.Fatalf("Reading uint32 failed. Wanted %X, got %X", expected, output)
	}
}

func TestReadUint64(t *testing.T) {
	expected := uint64(0x0102030405060708)
	raw := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	output := uint64(0)
	buffer := packetBuffer{bytes: raw}
	err := buffer.read(&output)
	if err != nil {
		t.Fatalf("Value parsing failed. Err = %v", err)
	}
	if output != expected {
		t.Fatalf("Reading uint64 failed. Wanted %X, got %X", expected, output)
	}
}

func TestReadStruct(t *testing.T) {
	type S struct {
		A uint8
		B uint16
		C uint32
		D uint64
	}
	expected := S{0x01, 0x0203, 0x04050607, 0x08090A0B0C0D0E0F}
	raw := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}
	var output S
	buffer := packetBuffer{bytes: raw}
	err := buffer.read(&output)
	if err != nil {
		t.Fatalf("Struct parsing failed. Err = %v", err)
	}
	if output != expected {
		t.Fatalf("Reading struct failed. Wanted %v, got %v", expected, output)
	}
}

func TestReadSlice(t *testing.T) {
	expected := []uint16{0x0102, 0x0304, 0x0506, 0x0708}
	raw := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	var output []uint16
	buffer := packetBuffer{bytes: raw}
	err := buffer.read(&output)
	if err != nil {
		t.Fatalf("Slice parsing failed. Err = %v", err)
	}
	if fmt.Sprintf("%x", output) != fmt.Sprintf("%x", expected) {
		t.Fatalf("Reading struct failed. Wanted %v, got %v", expected, output)
	}
}

func TestReadComplex(t *testing.T) {
	raw := []byte{
		0xf8,
		0x12, 0x34,
		0x56, 0x78, 0x9A, 0xBC,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x12,
		0x11, 0x22, 0x33, 0x44, 9,
		0x00, 0x00, 0x00, 0x01, 1, 0x00, 0x00, 0x00, 0x02, 5,
	}

	type Subtree struct {
		SubA uint32
		SubB uint8
	}

	type Tree struct {
		A uint8
		B uint16
		C uint32
		D uint64
		_ uint8
		F Subtree
		G []Subtree
	}

	expected := Tree{
		0xf8,
		0x1234,
		0x56789ABC,
		0x0102030405060708,
		0x00,
		Subtree{0x11223344, 9},
		[]Subtree{{0x01, 1}, {0x02, 5}},
	}

	var output Tree

	buffer := packetBuffer{bytes: raw}
	err := buffer.read(&output)
	if err != nil {
		t.Fatalf("Complex parsing failed. Err = %v", err)
	}
	if fmt.Sprintf("%x", output) != fmt.Sprintf("%x", expected) {
		t.Fatalf("Reading struct failed. Wanted %v, got %v", expected, output)
	}
}
