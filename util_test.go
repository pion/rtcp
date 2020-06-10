package rtcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPadding(t *testing.T) {
	assert := assert.New(t)
	type testCase struct {
		input  int
		result int
	}

	cases := []testCase{
		{input: 0, result: 0},
		{input: 1, result: 3},
		{input: 2, result: 2},
		{input: 3, result: 1},
		{input: 4, result: 0},
		{input: 100, result: 0},
		{input: 500, result: 0},
	}
	for _, testCase := range cases {
		assert.Equalf(getPadding(testCase.input), testCase.result, "Test case returned wrong value for input %d", testCase.input)
	}
}

func TestSetNBitsOfUint16(t *testing.T) {
	for _, test := range []struct {
		name   string
		source uint16
		size   uint16
		index  uint16
		value  uint16
		result uint16
		err    string
	}{
		{
			"setOneBit", 0, 1, 8, 1, 128, "",
		},
		{
			"setStatusVectorBit", 0, 1, 0, 1, 32768, "",
		},
		{
			"setStatusVectorSecondBit", 32768, 1, 1, 1, 49152, "",
		},
		{
			"setStatusVectorInnerBitsAndCutValue", 49152, 2, 6, 11111, 49920, "",
		},
		{
			"setRunLengthSecondTwoBit", 32768, 2, 1, 1, 40960, "",
		},
		{
			"setOneBitOutOfBounds", 32768, 2, 15, 1, 0, "invalid size or startIndex",
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := setNBitsOfUint16(test.source, test.size, test.index, test.value)
			if err != nil {
				if err.Error() != test.err {
					t.Fatalf("setNBitsOfUint16 %q : got = %v, want %v", test.name, err, test.err)
				}
				return
			}
			if got != test.result {
				t.Fatalf("setNBitsOfUint16 %q : got = %v, want %v", test.name, got, test.result)
			}
		})
	}
}
