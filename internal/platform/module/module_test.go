package module

import (
	"sync"
	"testing"
)

type testModule struct{ definition Definition }

func (m testModule) Definition() Definition { return m.definition }

func TestRegistryRejectsNilAndEmptyModules(t *testing.T) {
	r := NewRegistry()
	r.Register(nil)
	r.Register(testModule{})
	if got := r.List(); len(got) != 0 {
		t.Fatalf("expected empty registry, got %d modules", len(got))
	}
}

func TestRegistryListsDeterministically(t *testing.T) {
	r := NewRegistry()
	for _, name := range []string{"zeta", "alpha", "mesh"} {
		r.Register(testModule{definition: Definition{Name: name}})
	}
	got := r.List()
	want := []string{"alpha", "mesh", "zeta"}
	for i, name := range want {
		if got[i].Name != name {
			t.Fatalf("module %d: got %q, want %q", i, got[i].Name, name)
		}
	}
}

func TestRegistryConcurrentAccess(t *testing.T) {
	r := NewRegistry()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := "module-" + string(rune('a'+i))
			r.Register(testModule{definition: Definition{Name: name}})
			_, _ = r.Get(name)
			_ = r.Has(name)
			_ = r.List()
		}(i)
	}
	wg.Wait()
	if got := len(r.List()); got != 16 {
		t.Fatalf("expected 16 modules, got %d", got)
	}
}

func TestDependenciesReady(t *testing.T) {
	r := NewRegistry()
	r.Register(testModule{definition: Definition{Name: "dns"}})
	if !r.DependenciesReady(Definition{Dependencies: []string{"dns"}}) {
		t.Fatal("expected dependency to be ready")
	}
	if r.DependenciesReady(Definition{Dependencies: []string{"missing"}}) {
		t.Fatal("expected missing dependency to be unavailable")
	}
}
