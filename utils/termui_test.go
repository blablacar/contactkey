package utils

import "testing"

func TestVTClean(t *testing.T) {
	s := VTClean("[1;34m17:13:52[0m [1;33mWARN[0m[1;36m")
	if s != "17:13:52 WARN" {
		t.Errorf("Unexpected s : %q", s)
	}

	s = VTClean("[1A[K
	if s != "
		t.Errorf("Unexpected s : %q", s)
	}

}