package dns

import (
	"net"
	"strings"
)

// ReverseZone returns the standard reverse-DNS zone name for an IP address.
// The DNS adapter is responsible for publishing the resulting zone.
func ReverseZone(ip string) (string, bool) {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil { return "", false }
	if v4 := parsed.To4(); v4 != nil {
		return stringDigit(v4[3]) + "." + stringDigit(v4[2]) + "." + stringDigit(v4[1]) + "." + stringDigit(v4[0]) + ".in-addr.arpa.", true
	}
	b := parsed.To16()
	if b == nil { return "", false }
	const hex = "0123456789abcdef"
	var out strings.Builder
	for i := len(b) - 1; i >= 0; i-- {
		out.WriteByte(hex[b[i]&0x0f]); out.WriteByte('.')
		out.WriteByte(hex[b[i]>>4]); out.WriteByte('.')
	}
	out.WriteString("ip6.arpa.")
	return out.String(), true
}

func stringDigit(b byte) string {
	if b == 0 { return "0" }
	var buf [3]byte; i := len(buf)
	for b > 0 { i--; buf[i] = '0' + b%10; b /= 10 }
	return string(buf[i:])
}
