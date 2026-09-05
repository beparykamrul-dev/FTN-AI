package gis

import (
	"sort"
	"strings"
)

type ImpactNode struct {
	AssetID string `json:"asset_id"`
	Kind    string `json:"kind"`
	Name    string `json:"name,omitempty"`
	Status  string `json:"status,omitempty"`
}

type ImpactResult struct {
	Root     string       `json:"root"`
	Affected []ImpactNode `json:"affected"`
}

// ImpactAnalyzer traverses the modeled topology only. It does not execute
// network changes; remediation remains an explicit downstream action.
type ImpactAnalyzer struct {
	Graph *TopologyGraph
}

func NewImpactAnalyzer(g *TopologyGraph) *ImpactAnalyzer { return &ImpactAnalyzer{Graph: g} }

func (a *ImpactAnalyzer) Analyze(root string, nodes map[string]ImpactNode) ImpactResult {
	root = strings.TrimSpace(root)
	result := ImpactResult{Root: root, Affected: []ImpactNode{}}
	if a == nil || a.Graph == nil || root == "" || len(nodes) == 0 {
		return result
	}
	seen := map[string]bool{root: true}
	queue := []string{root}
	affected := make([]ImpactNode, 0)
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		if n, ok := nodes[id]; ok && id != root {
			affected = append(affected, n)
		}
		for _, l := range a.Graph.Neighbors(id) {
			next := l.ToID
			if next == id {
				next = l.FromID
			}
			if next != "" && !seen[next] {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	sort.SliceStable(affected, func(i, j int) bool {
		if affected[i].AssetID != affected[j].AssetID {
			return affected[i].AssetID < affected[j].AssetID
		}
		return affected[i].Kind < affected[j].Kind
	})
	result.Affected = affected
	return result
}
