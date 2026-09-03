package traffic

import (
	"sort"
	"strings"
	"time"
)

type Class string

const (
	ClassRealtime Class = "realtime"
	ClassVoice    Class = "voice"
	ClassVideo    Class = "video"
	ClassGaming   Class = "gaming"
	ClassNormal   Class = "normal"
	ClassBulk     Class = "bulk"
)

type Service struct {
	ID        string   `json:"id"`
	Class     Class    `json:"class"`
	Protocols []string `json:"protocols"`
	Priority  uint8    `json:"priority"`
	DSCP      uint8    `json:"dscp"`
}

// Registry intentionally uses stable service identities rather than hard-coded IPs.
// Provider endpoint changes are therefore handled by the classifier/feed layer.
func DefaultServices() []Service {
	return []Service{
		{ID: "whatsapp", Class: ClassRealtime, Protocols: []string{"udp", "tcp", "quic"}, Priority: 90, DSCP: 46},
		{ID: "telegram", Class: ClassRealtime, Protocols: []string{"udp", "tcp", "quic"}, Priority: 90, DSCP: 46},
		{ID: "imo", Class: ClassVoice, Protocols: []string{"udp", "tcp", "quic"}, Priority: 92, DSCP: 46},
		{ID: "pubg", Class: ClassGaming, Protocols: []string{"udp", "tcp"}, Priority: 95, DSCP: 46},
		{ID: "freefire", Class: ClassGaming, Protocols: []string{"udp", "tcp"}, Priority: 95, DSCP: 46},
		{ID: "realtime-generic", Class: ClassRealtime, Protocols: []string{"udp", "tcp", "quic", "webrtc"}, Priority: 85, DSCP: 46},
		{ID: "voice-generic", Class: ClassVoice, Protocols: []string{"udp", "quic", "webrtc"}, Priority: 88, DSCP: 46},
		{ID: "video-generic", Class: ClassVideo, Protocols: []string{"udp", "tcp", "quic", "webrtc"}, Priority: 75, DSCP: 34},
		{ID: "normal", Class: ClassNormal, Protocols: []string{"tcp", "udp"}, Priority: 40, DSCP: 0},
		{ID: "bulk", Class: ClassBulk, Protocols: []string{"tcp", "udp"}, Priority: 10, DSCP: 8},
	}
}

type Observation struct {
	PathID      string
	ServiceID   string
	Class       Class
	LatencyMs   float64
	JitterMs    float64
	PacketLoss  float64
	Congestion  float64
	Healthy     bool
	ObservedAt  time.Time
}

type Decision struct {
	PathID      string  `json:"path_id"`
	ServiceID   string  `json:"service_id"`
	Class       Class   `json:"class"`
	Score       float64 `json:"score"`
	DSCP        uint8   `json:"dscp"`
	Priority    uint8   `json:"priority"`
	Failover    bool    `json:"failover"`
	HoldDownSec int     `json:"hold_down_sec"`
}

func Score(o Observation) float64 {
	if !o.Healthy { return -1 }
	s := 100.0 - o.LatencyMs - (o.JitterMs * 1.5) - (o.PacketLoss * 5) - (o.Congestion * 20)
	if o.ObservedAt.IsZero() { return s - 10 }
	if age := time.Since(o.ObservedAt); age > 30*time.Second { s -= 20 }
	return s
}

func Select(observations []Observation, service Service) (Decision, bool) {
	best := Observation{}
	bestScore := -1.0
	for _, o := range observations {
		if strings.TrimSpace(o.ServiceID) != service.ID && strings.TrimSpace(o.Class) != string(service.Class) { continue }
		s := Score(o)
		if s > bestScore || (s == bestScore && o.PathID < best.PathID) { best, bestScore = o, s }
	}
	if bestScore < 0 || best.PathID == "" { return Decision{}, false }
	return Decision{PathID: best.PathID, ServiceID: service.ID, Class: service.Class, Score: bestScore, DSCP: service.DSCP, Priority: service.Priority, HoldDownSec: 5}, true
}

func SortServices(services []Service) []Service {
	out := append([]Service(nil), services...)
	sort.SliceStable(out, func(i, j int) bool { if out[i].Priority == out[j].Priority { return out[i].ID < out[j].ID }; return out[i].Priority > out[j].Priority })
	return out
}
