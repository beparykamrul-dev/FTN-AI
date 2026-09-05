package ftndns

import (
	"container/list"
	"strings"
	"sync"
	"time"
)

type CacheEntry struct { Key string; Value []string; ExpiresAt time.Time }
type ResolverCache struct { mu sync.RWMutex; maxEntries int; items map[string]*list.Element; lru *list.List }

func NewResolverCache(maxEntries int) *ResolverCache { if maxEntries < 1 { maxEntries = 1 }; return &ResolverCache{maxEntries: maxEntries, items: make(map[string]*list.Element), lru: list.New()} }
func (c *ResolverCache) Get(key string, now time.Time) ([]string, bool) {
	if c == nil { return nil, false }; key = strings.TrimSpace(key); if key == "" { return nil, false }
	c.mu.Lock(); defer c.mu.Unlock(); e, ok := c.items[key]; if !ok { return nil, false }; entry := e.Value.(CacheEntry)
	if !entry.ExpiresAt.After(now) { delete(c.items,key); c.lru.Remove(e); return nil,false }; c.lru.MoveToFront(e); return append([]string(nil),entry.Value...),true
}
func (c *ResolverCache) Set(key string, value []string, expiresAt time.Time) {
	if c == nil { return }; key = strings.TrimSpace(key); if key == "" || !expiresAt.After(time.Time{}) { return }
	c.mu.Lock(); defer c.mu.Unlock(); if c.items == nil { c.items = make(map[string]*list.Element) }; if c.lru == nil { c.lru = list.New() }
	if e,ok:=c.items[key]; ok { e.Value=CacheEntry{Key:key,Value:append([]string(nil),value...),ExpiresAt:expiresAt}; c.lru.MoveToFront(e); return }
	e:=c.lru.PushFront(CacheEntry{Key:key,Value:append([]string(nil),value...),ExpiresAt:expiresAt}); c.items[key]=e
	for len(c.items)>c.maxEntries { old:=c.lru.Back(); if old==nil { break }; delete(c.items,old.Value.(CacheEntry).Key); c.lru.Remove(old) }
}
func (c *ResolverCache) PurgeExpired(now time.Time) int {
	if c == nil { return 0 }; c.mu.Lock(); defer c.mu.Unlock(); removed:=0
	for e:=c.lru.Back(); e!=nil; { prev:=e.Prev(); entry:=e.Value.(CacheEntry); if !entry.ExpiresAt.After(now) { delete(c.items,entry.Key); c.lru.Remove(e); removed++ }; e=prev }; return removed
}
