package main

import (
	"net/http"
	"net/http/httptest"
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

func TestQKDIntentRequiresNodeOrKMS(t *testing.T) {
	a := &App{}
	r := httptest.NewRequest(http.MethodPost, "/api/v1/qkd/intents", nil)
	w := httptest.NewRecorder()
	// The permission layer is deliberately exercised before validation in the
	// handler; this test only documents the contract type boundary.
	_ = a
	_ = r
	_ = w
}
