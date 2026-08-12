package websocket

import (
	"context"
	"testing"
)

func TestFabricPublishesAndFiltersTopics(t *testing.T) {
	f := NewFabric()
	control, err := f.Subscribe("control", "server")
	if err != nil {
		t.Fatal(err)
	}
	monitoring, err := f.Subscribe("monitoring", "monitoring")
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Publish(context.Background(), Event{Type: "server.health.changed", Source: "agent"}); err != nil {
		t.Fatal(err)
	}

	select {
	case event := <-control.Events:
		if event.Sequence != 1 || event.Version != ProtocolVersion {
			t.Fatalf("unexpected event envelope: %+v", event)
		}
	case <-monitoring.Events:
		t.Fatal("monitoring subscriber received unrelated server event")
	default:
		t.Fatal("control subscriber did not receive event")
	}
}

func TestFabricRejectsPublishAfterClose(t *testing.T) {
	f := NewFabric()
	f.Close()
	if err := f.Publish(context.Background(), Event{Type: "test"}); err != ErrClosed {
		t.Fatalf("publish error = %v, want %v", err, ErrClosed)
	}
}
