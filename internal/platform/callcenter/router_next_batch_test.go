package callcenter

import "testing"

func TestCallRouterConstructor(t *testing.T) {
	r := NewRouter()
	if r == nil { t.Fatal("constructor returned nil") }
	if err := r.Register("event.test", nil); err == nil { t.Fatal("nil handler must be rejected") }
}
