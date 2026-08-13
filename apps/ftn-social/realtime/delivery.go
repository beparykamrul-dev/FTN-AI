package realtime

import "sync"

type DeliveryState string

const (
	Queued   DeliveryState = "queued"
	Delivered DeliveryState = "delivered"
	Failed   DeliveryState = "failed"
)

type Delivery struct {
	MessageID string
	UserID    string
	State     DeliveryState
}

type Tracker struct {
	mu    sync.RWMutex
	items map[string]Delivery
}

func NewTracker() *Tracker {
	return &Tracker{items: make(map[string]Delivery)}
}

func (t *Tracker) Set(delivery Delivery) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.items[delivery.MessageID+":"+delivery.UserID] = delivery
}

func (t *Tracker) Get(messageID, userID string) (Delivery, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	delivery, ok := t.items[messageID+":"+userID]
	return delivery, ok
}
