package pools

import (
	"testing"
)

func TestGetPool(t *testing.T) {
	// TODO: Figure out how to find this on its own
	SetJsonBase("../data/json")
	_, err := GetPool("technician")
	if err != nil {
		t.Error("Test failed with error:", err)
	}
}
