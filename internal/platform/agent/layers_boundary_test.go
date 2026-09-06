package agent

import "testing"

func TestLayerRegistryPrefersNativeLayer(t *testing.T) {
	r := NewLayerRegistry([]Layer{
		{ID: "remote", Priority: 1, Categories: []Category{"chat"}, Enabled: true},
		{ID: "native", Priority: 9, Categories: []Category{"chat"}, Enabled: true, Native: true},
	})
	got, err := r.Resolve(Category("chat"))
	if err != nil || got.ID != "native" {
		t.Fatalf("unexpected layer: %+v, err=%v", got, err)
	}
}
