package dns

import (
 "context"
 "testing"
)

func TestGoBGPAdapterFailClosed(t *testing.T) {
 if err := (*GoBGPAdapter)(nil).Validate(); err == nil { t.Fatal("nil GoBGP adapter accepted") }
 a := NewGoBGPAdapter(" 127.0.0.1:50051 ", true)
 if err := a.Validate(); err != nil { t.Fatalf("valid GoBGP adapter rejected: %v", err) }
 if a.Address != "127.0.0.1:50051" { t.Fatalf("normalized address=%q", a.Address) }
 if err := a.Publish(context.Background(), nil); err != nil { t.Fatalf("empty publish rejected: %v", err) }
 if err := a.Publish(nil, nil); err == nil { t.Fatal("nil context accepted") }
}
