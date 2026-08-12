package ws

import (
	"context"
	"testing"
	"time"
)

type testClient struct { id string; events []Event }
func (c *testClient) ID() string { return c.id }
func (c *testClient) Send(_ context.Context, e Event) error { c.events = append(c.events, e); return nil }

func TestFabricPublishAndTopicMatch(t *testing.T) {
	f := NewFabric(8, time.Second)
	a := &testClient{id: "a"}
	b := &testClient{id: "b"}
	if err := f.Register(a); err != nil { t.Fatal(err) }
	if err := f.Register(b); err != nil { t.Fatal(err) }
	if err := f.Subscribe("a", "sub-a", "/monitoring/*"); err != nil { t.Fatal(err) }
	if err := f.Subscribe("b", "sub-b", "/control"); err != nil { t.Fatal(err) }
	if err := f.Publish(context.Background(), Event{Type: "server.health.changed", Topic: "/monitoring/servers", Source: "agent"}); err != nil { t.Fatal(err) }
	if len(a.events) != 1 || len(b.events) != 0 { t.Fatalf("unexpected routing: a=%d b=%d", len(a.events), len(b.events)) }
	if a.events[0].ID != 1 || a.events[0].Protocol != Protocol { t.Fatalf("bad event envelope: %+v", a.events[0]) }
}

func TestFabricUnregisterAndClose(t *testing.T) {
	f := NewFabric(1, time.Second)
	c := &testClient{id: "c"}
	if err := f.Register(c); err != nil { t.Fatal(err) }
	f.Unregister("c")
	if f.ClientCount() != 0 { t.Fatalf("client count = %d", f.ClientCount()) }
	f.Close()
	if err := f.Register(&testClient{id: "d"}); err != ErrClosed { t.Fatalf("register after close = %v", err) }
}
