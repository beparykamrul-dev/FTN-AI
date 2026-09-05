package gis

import (
	"sort"
	"strings"
	"sync"
)

type MapStore struct { mu sync.RWMutex; nodes map[string]MapNode; edges map[string]MapEdge }
func NewMapStore() *MapStore { return &MapStore{nodes: make(map[string]MapNode), edges: make(map[string]MapEdge)} }
func (m *MapStore) UpsertNode(n MapNode) { if m == nil { return }; n.ID=strings.TrimSpace(n.ID); if n.ID==""{return}; m.mu.Lock(); defer m.mu.Unlock(); if m.nodes==nil{m.nodes=make(map[string]MapNode)}; m.nodes[n.ID]=n }
func (m *MapStore) UpsertEdge(e MapEdge) { if m == nil { return }; e.ID=strings.TrimSpace(e.ID); if e.ID==""{return}; m.mu.Lock(); defer m.mu.Unlock(); if m.edges==nil{m.edges=make(map[string]MapEdge)}; m.edges[e.ID]=e }
func (m *MapStore) Snapshot() MapSnapshot { if m == nil { return MapSnapshot{} }; m.mu.RLock(); defer m.mu.RUnlock(); s:=MapSnapshot{Nodes:make([]MapNode,0,len(m.nodes)),Edges:make([]MapEdge,0,len(m.edges))}; for _,n:=range m.nodes{s.Nodes=append(s.Nodes,n)}; for _,e:=range m.edges{s.Edges=append(s.Edges,e)}; sort.Slice(s.Nodes,func(i,j int)bool{return s.Nodes[i].ID<s.Nodes[j].ID}); sort.Slice(s.Edges,func(i,j int)bool{return s.Edges[i].ID<s.Edges[j].ID}); return s }
