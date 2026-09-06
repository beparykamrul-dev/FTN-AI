package edge

import (
	"net"
	"testing"
)

func TestOLTStateRejectsInvalidManagementIP(t *testing.T) {
	if net.ParseIP("not-an-ip") != nil { t.Fatal("fixture unexpectedly parsed") }
	if (OLTState{ID:"olt-1", Vendor:"v", ManagementIP:"not-an-ip"}).Valid() { t.Fatal("invalid management IP accepted") }
}
