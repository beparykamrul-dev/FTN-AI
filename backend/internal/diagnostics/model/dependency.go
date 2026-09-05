package model

import "strings"

// Dependency represents an authorized relationship in the FTN service graph.
type Dependency struct { From string `json:"from"`; To string `json:"to"`; Kind string `json:"kind"`; Critical bool `json:"critical"`; Healthy bool `json:"healthy"` }

func (d Dependency) Valid() bool { return strings.TrimSpace(d.From) != "" && strings.TrimSpace(d.To) != "" && strings.TrimSpace(d.From) != strings.TrimSpace(d.To) }

func (d Dependency) Normalize() Dependency { d.From = strings.TrimSpace(d.From); d.To = strings.TrimSpace(d.To); d.Kind = strings.TrimSpace(d.Kind); return d }
