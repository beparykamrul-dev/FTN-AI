package realtime

import (
	"strings"
	"sync"
)

type DeliveryState string

const (
	Queued DeliveryState = "queued"
	Delivered DeliveryState = "delivered"
	Failed DeliveryState = "failed"
)

type Delivery struct { MessageID string; UserID string; State DeliveryState }
type Tracker struct { mu sync.RWMutex; items map[string]Delivery }
func NewTracker() *Tracker { return &Tracker{items: make(map[string]Delivery)} }
func normalizeDelivery(d Delivery) Delivery { d.MessageID=strings.TrimSpace(d.MessageID); d.UserID=strings.TrimSpace(d.UserID); return d }
func validDelivery(d Delivery) bool { d=normalizeDelivery(d); if d.MessageID==""||d.UserID==""||len(d.MessageID)>256||len(d.UserID)>256{return false}; switch d.State{case Queued,Delivered,Failed:return true};return false }
func deliveryKey(messageID,userID string) string { return messageID+":"+userID }
func (t *Tracker) Set(d Delivery) { if t==nil{return};d=normalizeDelivery(d);if !validDelivery(d){return};t.mu.Lock();defer t.mu.Unlock();if t.items==nil{t.items=make(map[string]Delivery)};t.items[deliveryKey(d.MessageID,d.UserID)]=d }
func (t *Tracker) Get(messageID,userID string)(Delivery,bool){if t==nil{return Delivery{},false};messageID=strings.TrimSpace(messageID);userID=strings.TrimSpace(userID);if messageID==""||userID==""||len(messageID)>256||len(userID)>256{return Delivery{},false};t.mu.RLock();defer t.mu.RUnlock();d,ok:=t.items[deliveryKey(messageID,userID)];return d,ok}
