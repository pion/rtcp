package rtcp

import "testing"

func TestUnmarshalNil(t *testing.T) {
	_, err := Unmarshal(nil)
	if got, want := err, errInvalidHeader; got != want {
		t.Fatalf("Unmarshal(nil) err = %v, want %v", got, want)
	}
}
