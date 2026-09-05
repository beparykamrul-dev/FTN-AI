package mesh

import "testing"

func TestDefaultDeviceProfileCanJoin(t *testing.T) {
	p := DefaultDeviceProfile("router")
	if !p.CanJoin() || p.KeepaliveSeconds == 0 { t.Fatalf("unsafe default device profile: %#v", p) }
}

func TestDeviceProfileCannotJoinWithoutMesh(t *testing.T) {
	p := DefaultDeviceProfile("router"); p.MeshEnabled = false
	if p.CanJoin() { t.Fatal("mesh-disabled profile must not join") }
}
