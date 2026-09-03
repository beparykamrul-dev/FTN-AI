package ftnftp

import (
	"fmt"
	"sync"
)

type Manager struct { mu sync.RWMutex; objects map[string]Object }
func NewManager()*Manager{return &Manager{objects:map[string]Object{}}}
func(m *Manager) Put(o Object)error{if o.ID==""||o.Bucket==""||o.Key==""{return fmt.Errorf("id, bucket and key are required")};m.mu.Lock();defer m.mu.Unlock();m.objects[o.ID]=o;return nil}
func(m *Manager) Get(id string)(Object,bool){m.mu.RLock();defer m.mu.RUnlock();o,ok:=m.objects[id];return o,ok}
func(m *Manager) List()[]Object{m.mu.RLock();defer m.mu.RUnlock();out:=make([]Object,0,len(m.objects));for _,o:=range m.objects{out=append(out,o)};return out}
