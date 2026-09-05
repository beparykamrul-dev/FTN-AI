package realtime

import "sync"

type DeliveryState string

const (
	Queued    DeliveryState = "queued"
	Delivered DeliveryState = "delivered"
	Failed    DeliveryState = "failed"
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

func validDelivery(d Delivery) bool {
	if d.MessageID == "" || d.UserID == "" || len(d.MessageID) > 256 || len(d.UserID) > 256 {
		return false
	}
	switch d.State {
	case Queued, Delivered, Failed:
		return true
	}
	return false
}

func (t *Tracker) Set(d Delivery) {
	if t == nil || !validDelivery(d) {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.items == nil {
		t.items = make(map[string]Delivery)
	}
	t.items[d.MessageID+":"+d.UserID] = d
}

func (t *Tracker) Get(messageID, userID string) (Delivery, bool) {
	if t == nil || messageID == "" || userID == "" || len(messageID) > 256 || len(userID) > 256 {
		return Delivery{}, false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	d, ok := t.items[messageID+":"+userID]
	return d, ok
}
