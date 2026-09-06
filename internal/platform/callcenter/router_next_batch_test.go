package callcenter

import "testing"

func TestCallRouterConstructorAndNilHandler(t *testing.T) {
	r := NewRouter()
	if r == nil { t.Fatal("constructor returned nil") }
	r.Register("event.test", nil)
	if errs := r.Dispatch(CallEvent{Type:""}); len(errs) != 1 { t.Fatalf("errs=%d, want one validation error", len(errs)) }
}
