package db

import (
	"context"
	"testing"
)

func TestRepositoryNilPoolIsSafe(t *testing.T) {
	r := &Repository{}
	if err := r.SaveServiceRequest(context.Background(), ServiceRequest{ServiceID:"billing"}); err != nil { t.Fatalf("unexpected error: %v", err) }
}
