package builder

import (
	"fmt"
	"os"
	"path/filepath"
)

// Generator creates a safe project skeleton from an approved manifest.
// It only writes below the supplied destination and never executes generated code.
type Generator struct{}

func NewGenerator() *Generator { return &Generator{} }

func (g *Generator) Generate(m Manifest, destination string) error {
	if err := m.Validate(); err != nil { return err }
	if destination == "" { return fmt.Errorf("destination is required") }

	root, err := filepath.Abs(destination)
	if err != nil { return err }
	if err := os.MkdirAll(root, 0o755); err != nil { return err }

	t, err := GetTemplate(m.Template)
	if err != nil { return err }

	for _, name := range t.Files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := ensureWithin(root, path); err != nil { return err }
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return err }
		if _, err := os.Stat(path); err == nil { continue }
		content := []byte("// FTN generated project skeleton: " + name + "\n")
		if err := os.WriteFile(path, content, 0o644); err != nil { return err }
	}
	return nil
}

func ensureWithin(root, path string) error {
	rel, err := filepath.Rel(root, path)
	if err != nil { return err }
	if rel == ".." || len(rel) >= 3 && rel[:3] == ".."+string(filepath.Separator) {
		return fmt.Errorf("path escapes project destination")
	}
	return nil
}
