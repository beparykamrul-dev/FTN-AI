package dns

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

type DNSQueryProbe struct {
	Address string        `json:"address"`
	Name    string        `json:"name"`
	Timeout time.Duration `json:"timeout"`
}

type DNSQueryProbeResult struct {
	Reachable bool          `json:"reachable"`
	Answered  bool          `json:"answered"`
	Latency   time.Duration `json:"latency"`
	RCode     uint8         `json:"rcode"`
	Error     string        `json:"error,omitempty"`
}

func (p DNSQueryProbe) Validate() error {
	if strings.TrimSpace(p.Address) == "" {
		return fmt.Errorf("DNS server address is required")
	}
	name := strings.TrimSuffix(strings.TrimSpace(p.Name), ".")
	if name == "" {
		return fmt.Errorf("DNS query name is required")
	}
	if len(name) > 253 || p.Timeout <= 0 {
		return fmt.Errorf("invalid DNS query name or timeout")
	}
	for _, label := range strings.Split(name, ".") {
		if len(label) == 0 || len(label) > 63 {
			return fmt.Errorf("invalid DNS query name")
		}
	}
	return nil
}

func (p DNSQueryProbe) Probe(ctx context.Context) DNSQueryProbeResult {
	result := DNSQueryProbeResult{}
	if ctx == nil {
		result.Error = "context is required"
		return result
	}
	if err := p.Validate(); err != nil {
		result.Error = err.Error()
		return result
	}
	if err := ctx.Err(); err != nil {
		result.Error = err.Error()
		return result
	}
	server := strings.TrimSpace(p.Address)
	if _, _, err := net.SplitHostPort(server); err != nil {
		server = net.JoinHostPort(server, "53")
	}
	id := uint16(time.Now().UnixNano())
	name := strings.TrimSuffix(strings.TrimSpace(p.Name), ".")
	packet := make([]byte, 12, 64)
	binary.BigEndian.PutUint16(packet[0:2], id)
	binary.BigEndian.PutUint16(packet[2:4], 0x0100)
	binary.BigEndian.PutUint16(packet[4:6], 1)
	for _, label := range strings.Split(name, ".") {
		packet = append(packet, byte(len(label)))
		packet = append(packet, label...)
	}
	packet = append(packet, 0, 0, 1, 0, 1)
	d := net.Dialer{Timeout: p.Timeout}
	conn, err := d.DialContext(ctx, "udp", server)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer conn.Close()
	deadline := time.Now().Add(p.Timeout)
	if dl, ok := ctx.Deadline(); ok && dl.Before(deadline) {
		deadline = dl
	}
	if err := conn.SetDeadline(deadline); err != nil {
		result.Error = err.Error()
		return result
	}
	start := time.Now()
	if _, err = conn.Write(packet); err != nil {
		result.Error = err.Error()
		return result
	}
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	result.Latency = time.Since(start)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	result.Reachable = true
	if n < 12 || binary.BigEndian.Uint16(buf[0:2]) != id {
		result.Error = "invalid DNS response"
		return result
	}
	result.Answered = binary.BigEndian.Uint16(buf[6:8]) > 0
	result.RCode = buf[3] & 0x0f
	return result
}
