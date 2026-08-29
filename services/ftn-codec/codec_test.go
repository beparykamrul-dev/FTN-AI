package codec

import "testing"

func TestDefaultCapabilities(t *testing.T) {
	caps := DefaultCapabilities()
	if len(caps) < 6 {
		t.Fatalf("capabilities=%d, want at least 6", len(caps))
	}
	seen := map[string]bool{}
	for _, c := range caps {
		if c.ID == "" || c.Class == "" {
			t.Fatalf("invalid capability: %+v", c)
		}
		if seen[c.ID] {
			t.Fatalf("duplicate capability %q", c.ID)
		}
		seen[c.ID] = true
	}
	if !seen["hardware-video-encode"] || !seen["chunking"] || !seen["deduplication"] {
		t.Fatal("required media/transfer capabilities are missing")
	}
}

func TestMediaWorkersAreIsolated(t *testing.T) {
	for _, c := range DefaultCapabilities() {
		if c.Class == "media" && !c.WorkerIsolated {
			t.Fatalf("media capability %q must be worker-isolated", c.ID)
		}
	}
}
