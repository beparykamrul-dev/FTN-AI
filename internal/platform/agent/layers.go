package agent

import (
	"fmt"
	"sort"
)

// Layer describes an approved intelligence layer. FTN-native layers are preferred.
type Layer struct {
	ID         string
	Priority   int
	Categories []Category
	Enabled    bool
	Native     bool
}

// LayerRegistry holds the centrally managed approved layers.
type LayerRegistry struct {
	layers map[string]Layer
}

func NewLayerRegistry(layers []Layer) *LayerRegistry {
	r := &LayerRegistry{layers: make(map[string]Layer, len(layers))}
	for _, l := range layers {
		if l.ID != "" {
			r.layers[l.ID] = l
		}
	}
	return r
}

func (r *LayerRegistry) Resolve(category Category) (Layer, error) {
	if r == nil {
		return Layer{}, fmt.Errorf("layer registry is required")
	}
	var candidates []Layer
	for _, l := range r.layers {
		if !l.Enabled {
			continue
		}
		for _, c := range l.Categories {
			if c == category {
				candidates = append(candidates, l)
				break
			}
		}
	}
	if len(candidates) == 0 {
		return Layer{}, fmt.Errorf("no enabled layer for category: %s", category)
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Native != candidates[j].Native {
			return candidates[i].Native
		}
		if candidates[i].Priority != candidates[j].Priority {
			return candidates[i].Priority < candidates[j].Priority
		}
		return candidates[i].ID < candidates[j].ID
	})
	return candidates[0], nil
}
