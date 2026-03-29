// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrPacketTooShortForPacket(t *testing.T) {
	t.Parallel()

	err := errPacketTooShortFor(&PictureLossIndication{})

	assert.ErrorIs(t, err, errPacketTooShort)
	assert.Contains(t, err.Error(), "PictureLossIndication")
}

func TestErrPacketTooShortForString(t *testing.T) {
	t.Parallel()

	err := errPacketTooShortFor("CustomPacket")

	assert.ErrorIs(t, err, errPacketTooShort)
	assert.Contains(t, err.Error(), "CustomPacket")
}

func TestErrPacketTooShortForNil(t *testing.T) {
	t.Parallel()

	err := errPacketTooShortFor(nil)

	assert.True(t, errors.Is(err, errPacketTooShort))
	assert.Equal(t, errPacketTooShort, err)
}

func TestPacketNameFromHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header Header
		want   string
	}{
		{
			name: "sender report",
			header: Header{
				Type: TypeSenderReport,
			},
			want: "SenderReport",
		},
		{
			name: "transport feedback",
			header: Header{
				Type:  TypeTransportSpecificFeedback,
				Count: FormatTLN,
			},
			want: "TransportLayerNack",
		},
		{
			name: "transport feedback fallback",
			header: Header{
				Type:  TypeTransportSpecificFeedback,
				Count: 99,
			},
			want: "TransportSpecificFeedback(FMT=99)",
		},
		{
			name: "payload specific fallback",
			header: Header{
				Type:  TypePayloadSpecificFeedback,
				Count: 5,
			},
			want: "PayloadSpecificFeedback(FMT=5)",
		},
		{
			name: "payload specific known",
			header: Header{
				Type:  TypePayloadSpecificFeedback,
				Count: FormatPLI,
			},
			want: "PictureLossIndication",
		},
		{
			name: "unknown type",
			header: Header{
				Type: 199,
			},
			want: "PacketType(199)",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, packetNameFromHeader(test.header))
		})
	}
}
