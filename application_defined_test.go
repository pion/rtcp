// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTApplicationPacketUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Data      []byte
		Want      ApplicationDefined
		WantError error
	}{
		{
			Name: "valid",
			Data: []byte{
				// Application Packet Type + Length(0x0003)
				0x80, 0xcc, 0x00, 0x03,
				// sender=0x4baae1ab
				0x4b, 0xaa, 0xe1, 0xab,
				// name='NAME'
				0x4E, 0x41, 0x4D, 0x45,
				// data='ABCD'
				0x41, 0x42, 0x43, 0x44,
			},
			Want: ApplicationDefined{
				SubType: 0,
				SSRC:    0x4baae1ab,
				Name:    "NAME",
				Data:    []byte{0x41, 0x42, 0x43, 0x44},
			},
		},
		{
			Name: "validCustomSsubType",
			Data: []byte{
				// Application Packet Type (SubType 31) + Length(0x0003)
				0x9f, 0xcc, 0x00, 0x03,
				// sender=0x4baae1ab
				0x4b, 0xaa, 0xe1, 0xab,
				// name='NAME'
				0x4E, 0x41, 0x4D, 0x45,
				// data='ABCD'
				0x41, 0x42, 0x43, 0x44,
			},
			Want: ApplicationDefined{
				SubType: 31,
				SSRC:    0x4baae1ab,
				Name:    "NAME",
				Data:    []byte{0x41, 0x42, 0x43, 0x44},
			},
		},
		{
			Name: "validWithPadding",
			Data: []byte{
				// Application Packet Type + Length(0x0002)  (0xA0 has padding bit set)
				0xA0, 0xcc, 0x00, 0x04,
				// sender=0x4baae1ab
				0x4b, 0xaa, 0xe1, 0xab,
				// name='NAME'
				0x4E, 0x41, 0x4D, 0x45,
				// data='ABCDE'
				0x41, 0x42, 0x43, 0x44, 0x45,
				// 3 bytes padding as packet length must be a division of 4
				0x03, 0x03, 0x03,
			},
			Want: ApplicationDefined{
				SubType: 0,
				SSRC:    0x4baae1ab,
				Name:    "NAME",
				Data:    []byte{0x41, 0x42, 0x43, 0x44, 0x45},
			},
		},
		{
			Name: "invalidAppPacketLengthField",
			Data: []byte{
				// Application Packet Type + invalid Length(0x00FF)
				0x80, 0xcc, 0x00, 0xFF,
				// sender=0x4baae1ab
				0x4b, 0xaa, 0xe1, 0xab,
				// name='NAME'
				0x4E, 0x41, 0x4D, 0x45,
				// data='ABCD'
				0x41, 0x42, 0x43, 0x44,
			},
			WantError: errAppDefinedInvalidLength,
		},
		{
			Name: "invalidPacketLengthTooShort",
			Data: []byte{
				// Application Packet Type + Length(0x0002). Total packet length is less than 12 bytes
				0x80, 0xcc, 0x00, 0x2,
				// sender=0x4baae1ab
				0x4b, 0xaa, 0xe1, 0xab,
				// name='SUI'
				0x53, 0x55, 0x49,
			},
			WantError: errPacketTooShort,
		},
		{
			Name: "wrongPaddingSize",
			Data: []byte{
				// Application Packet Type + Length(0x0002)  (0xA0 has padding bit set)
				0xA0, 0xcc, 0x00, 0x04,
				// sender=0x4baae1ab
				0x4b, 0xaa, 0xe1, 0xab,
				// name='NAME'
				0x4E, 0x41, 0x4D, 0x45,
				// data='ABCDE'
				0x41, 0x42, 0x43, 0x44, 0x45,
				// 3 bytes padding as packet length must be a division of 4
				0x03, 0x03, 0x09, // last byte has padding size 0x09 which is more than the data + padding bytes
			},
			WantError: errWrongPadding,
		},
		{
			Name: "invalidHeader",
			Data: []byte{
				// Application Packet Type + invalid Length(0x00FF)
				0xFF,
			},
			WantError: errPacketTooShort,
		},
	} {
		var apk ApplicationDefined
		err := apk.Unmarshal(test.Data)
		assert.ErrorIsf(t, err, test.WantError, "Unmarshal %q", test.Name)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.Want, apk, "Unmarshal %q", test.Name)

		// Check SSRC is matching
		assert.Equalf(t, uint32(0x4baae1ab), apk.SSRC, "%q SSRC mismatch", test.Name)
		assert.Equalf(t, uint32(0x4baae1ab), apk.DestinationSSRC()[0], "%q DestinationSSRC mismatch", test.Name)
	}
}

func TestTApplicationPacketMarshal(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Want      []byte
		Packet    ApplicationDefined
		WantError error
	}{
		{
			Name: "valid",
			Want: []byte{
				// Application Packet Type + Length(0x0003)
				0x80, 0xcc, 0x00, 0x03,
				// sender=0x4baae1ab
				0x4b, 0xaa, 0xe1, 0xab,
				// name='NAME'
				0x4E, 0x41, 0x4D, 0x45,
				// data='ABCD'
				0x41, 0x42, 0x43, 0x44,
			},
			Packet: ApplicationDefined{
				SSRC: 0x4baae1ab,
				Name: "NAME",
				Data: []byte{0x41, 0x42, 0x43, 0x44},
			},
		},
		{
			Name: "validCustomSubType",
			Want: []byte{
				// Application Packet Type (SubType 31) + Length(0x0003)
				0x9f, 0xcc, 0x00, 0x03,
				// sender=0x4baae1ab
				0x4b, 0xaa, 0xe1, 0xab,
				// name='NAME'
				0x4E, 0x41, 0x4D, 0x45,
				// data='ABCD'
				0x41, 0x42, 0x43, 0x44,
			},
			Packet: ApplicationDefined{
				SubType: 31,
				SSRC:    0x4baae1ab,
				Name:    "NAME",
				Data:    []byte{0x41, 0x42, 0x43, 0x44},
			},
		},
		{
			Name: "validWithPadding",
			Want: []byte{
				// Application Packet Type + Length(0x0002)  (0xA0 has padding bit set)
				0xA0, 0xcc, 0x00, 0x04,
				// sender=0x4baae1ab
				0x4b, 0xaa, 0xe1, 0xab,
				// name='NAME'
				0x4E, 0x41, 0x4D, 0x45,
				// data='ABCDE'
				0x41, 0x42, 0x43, 0x44, 0x45,
				// 3 bytes padding as packet length must be a division of 4
				0x03, 0x03, 0x03,
			},
			Packet: ApplicationDefined{
				SSRC: 0x4baae1ab,
				Name: "NAME",
				Data: []byte{0x41, 0x42, 0x43, 0x44, 0x45},
			},
		},
		{
			Name:      "invalidDataTooLarge",
			WantError: errAppDefinedDataTooLarge,
			Packet: ApplicationDefined{
				SSRC: 0x4baae1ab,
				Name: "NAME",
				Data: make([]byte, 0xFFFF-12+1), // total max packet size is 0xFFFF including header and other fields.
			},
		},
		{
			Name:      "invalidName",
			WantError: errAppDefinedInvalidName,
			Packet: ApplicationDefined{
				SSRC: 0x4baae1ab,
				Name: "NOT4CHARS",
				Data: []byte{0x41, 0x42, 0x43, 0x44},
			},
		},
		{
			Name:      "InvalidSubType",
			WantError: errInvalidHeader,
			Packet: ApplicationDefined{
				SubType: 32, // Must be up to 31
				SSRC:    0x4baae1ab,
				Name:    "NAME",
				Data:    []byte{0x41, 0x42, 0x43, 0x44},
			},
		},
	} {
		rawPacket, err := test.Packet.Marshal()
		assert.ErrorIsf(t, err, test.WantError, "Marshal %q", test.Name)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.Want, rawPacket, "Marshal %q", test.Name)

		marshalSize := test.Packet.MarshalSize()
		assert.Equalf(t, marshalSize, len(rawPacket), "MarshalSize %q", test.Name)
	}
}
