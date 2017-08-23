package utils

import "testing"

func TestVTClean(t *testing.T) {
	s := VTClean("[1;34m17:13:52[0m [1;33mWARN[0m[1;36m")
	if s != "17:13:52 WARN" {
		t.Errorf("Unexpected s : %q", s)
	}

	s = VTClean("[1A[K[0;35m[2m[ZkCheck][webhooks][22m [1mwebhooks2[21m[22m /services/webhooks - Checking that service adds key in zookeeper (timeout in 595s)[0m")
	if s != "[ZkCheck][webhooks] webhooks2 /services/webhooks - Checking that service adds key in zookeeper (timeout in 595s)" {
		t.Errorf("Unexpected s : %q", s)
	}

}
