package main

import (
	"encoding/json"
	"testing"
)

func TestQKDNodeDoesNotExposeRawKeyMaterial(t *testing.T) {
	if _, ok := any(QKDNode{}).(interface{ RawKey() }); ok {
		t.Fatal("QKDNode must not expose raw key material")
	}
	if _, ok := any(QKDStatus{}).(interface{ RawKey() }); ok {
		t.Fatal("QKDStatus must not expose raw key material")
	}
}

func TestQKDContractsSerializeOnlyMetadata(t *testing.T) {
	node := QKDNode{ID: "node-1", Name: "qkd-1", Role: "kms", KMSRef: "kms://ftn/qkd-1", Status: "healthy"}
	status := QKDStatus{NodeID: node.ID, KMSRef: node.KMSRef, PoolState: "ready", AvailableKeys: 10, Healthy: true}
	for name, value := range map[string]any{"node": node, "status": status} {
		body, err := json.Marshal(value)
		if err != nil {
			t.Fatalf("%s marshal failed: %v", name, err)
		}
		if len(body) == 0 || containsRawKeyMaterial(body) {
			t.Fatalf("%s contract exposed unexpected key material: %s", name, body)
		}
	}
}

func containsRawKeyMaterial(body []byte) bool {
	var decoded map[string]any
	if json.Unmarshal(body, &decoded) != nil {
		return true
	}
	for _, key := range []string{"raw_key", "raw_key_material", "key_material", "secret_key", "private_key"} {
		if _, ok := decoded[key]; ok {
			return true
		}
	}
	return false
}
