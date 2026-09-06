package codec

import "testing"

func TestCodecRejectsNegativeResultBytes(t *testing.T) {
	if (Result{JobID:"j", Status:"s", BytesIn:-1}).Valid() { t.Fatal("negative input bytes accepted") }
	if (Job{CapabilityID:"c", InputURI:"in", Options:map[string]string{"":"x"}}).Valid() == nil { t.Fatal("empty option key accepted") }
}
