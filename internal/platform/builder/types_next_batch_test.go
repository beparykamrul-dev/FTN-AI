package builder

import "testing"

func TestProjectAndBuildJobValidation(t *testing.T) {
	p := Project{ID:"p1", Name:"demo", Target:TargetWeb, Backend:BackendSpec{Language:"go", Framework:"net/http"}, Frontend:FrontendSpec{Framework:"react", Platform:"web"}}
	if !p.Valid() { t.Fatal("valid project rejected") }
	p.Target = "invalid"
	if p.Valid() { t.Fatal("invalid target accepted") }
	if (BuildJob{ID:"j1", ProjectID:"p1", Status:"unknown"}).Valid() { t.Fatal("unknown build status accepted") }
}
