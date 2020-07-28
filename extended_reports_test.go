package rtcp

import (
	"reflect"
	"testing"
)

func TestExtendedReportsUnmarshal(t *testing.T) {
	for _, test := range []struct {
		Name string
		Data []byte
		Want ExtendedReports
		WantError error
	}{
		Name: "nil",
		Data: nil,
		WantError: errPacketTooShort,
	},
	{	
		Name: "valid",
		Data: []byte{

		}


	}
}

func makeExtendedReportsHeader() []byte {

}

func makeXRBT1() []byte {
	var report := &XRBT1{
		BlockType: 8,
		Thinning: 3,

	}
}

func makeXRBT2() []byte {

}

func makeXRBT3() []byte {

}

func makeXRBT4() []byte {

}

func makeXRBT5() []byte {

}

func makeXRBT6() []byte {

}

func makeXRBT7() []byte {

}

func makeXRBT8() []byte {

}

func makeXRBT10() []byte {

}

func makeXRBT11() []byte {

}

func makeXRBT12() []byte {

}

func makeXRBT13() []byte {

}

func makeXRBT14() []byte {

}

func makeXRBT15() []byte {

}

func makeXRBT16() []byte {

}

func makeXRBT17() []byte {

}

func makeXRBT18() []byte {

}

func makeXRBT19() []byte {

}

func makeXRBT20() []byte {

}

func makeXRBT21() []byte {

}

func makeXRBT22() []byte {

}

func makeXRBT23() []byte {

}

func makeXRBT24() []byte {

}

func makeXRBT25() []byte {

}

func makeXRBT26() []byte {

}

func makeXRBT27() []byte {

}

func makeXRBT28() []byte {

}

func makeXRBT29() []byte {

}

func makeXRBT30() []byte {

}

func makeXRBT31() []byte {

}

func makeXRBT32() []byte {

}

func makeXRBT33() []byte {

}

func makeXRBT34() []byte {

}

func makeXRBT35() []byte {

}

func changeIntToHex() {

}