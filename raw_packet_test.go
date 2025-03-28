// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Packet = (*RawPacket)(nil) // assert is a Packet

func TestRawPacketRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name               string
		Packet             RawPacket
		WantMarshalError   error
		WantUnmarshalError error
	}{
		{
			Name: "valid",
			Packet: RawPacket([]byte{
				// v=2, p=0, count=1, BYE, len=12
				0x81, 0xcb, 0x00, 0x0c,
				// ssrc=0x902f9e2e
				0x90, 0x2f, 0x9e, 0x2e,
				// len=3, text=FOO
				0x03, 0x46, 0x4f, 0x4f,
			}),
		},
		{
			Name:               "short header",
			Packet:             RawPacket([]byte{0x00}),
			WantUnmarshalError: errPacketTooShort,
		},
		{
			Name: "invalid header",
			Packet: RawPacket([]byte{
				// v=0, p=0, count=0, RR, len=4
				0x00, 0xc9, 0x00, 0x04,
			}),
			WantUnmarshalError: errBadVersion,
		},
	} {
		data, err := test.Packet.Marshal()
		assert.ErrorIsf(t, err, test.WantMarshalError, "Marshal %q", test.Name)
		if err != nil {
			continue
		}

		var decoded RawPacket

		err = decoded.Unmarshal(data)
		assert.ErrorIsf(t, err, test.WantUnmarshalError, "Unmarshal %q", test.Name)
		if err != nil {
			continue
		}
		assert.Equalf(t, test.Packet, decoded, "Unmarshal %q", test.Name)
	}
}
