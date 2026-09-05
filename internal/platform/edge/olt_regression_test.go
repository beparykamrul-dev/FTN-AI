package edge

import (
	"net"
	"testing"
)

func TestOLTStateRejectsInvalidManagementIP(t *testing.T) {
	ip := net.ParseIP("not-an-ip")
	if ip != nil { t.Fatal("test fixture unexpectedly parsed") }
	if (OLTState{ID:"olt-1", Vendor:"v", ManagementIP:"not-an-ip"}).Valid() { t.Fatal("invalid management IP must be rejected") }
}
