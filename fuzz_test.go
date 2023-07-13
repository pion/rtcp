// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtcp

import (
	"testing"
)

func FuzzUnmarshal(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		packets, err := Unmarshal(data)
		if err != nil {
			return
		}

		for _, packet := range packets {
			packet.Marshal()
		}
	})
}
