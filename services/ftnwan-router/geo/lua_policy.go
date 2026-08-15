package geo

import (
	"errors"
	"strings"
)

// LuaPolicy describes a sandboxed Lua-compatible policy expression.
// Execution is intentionally delegated to a future sandbox runtime; this
// layer only validates and stores policy metadata.
type LuaPolicy struct {
	Name       string
	Expression string
	Enabled    bool
}

var ErrInvalidPolicy = errors.New("invalid geo policy")

// ValidatePolicy performs safe structural validation before a policy can be
// handed to an isolated Lua runtime. It does not execute arbitrary code.
func ValidatePolicy(p LuaPolicy) error {
	if strings.TrimSpace(p.Name) == "" || strings.TrimSpace(p.Expression) == "" {
		return ErrInvalidPolicy
	}
	if len(p.Name) > 128 || len(p.Expression) > 4096 {
		return ErrInvalidPolicy
	}
	return nil
}
