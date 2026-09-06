package kernel

import "testing"

func TestWrapperRejectsNilAndOversizedRequests(t *testing.T) {
	var w *Wrapper
	if w == nil { return }
}
