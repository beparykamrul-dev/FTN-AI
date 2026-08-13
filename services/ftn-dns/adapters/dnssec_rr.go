package adapters

import (
	"encoding/binary"
	"fmt"
)

// DNSSECRR contains normalized metadata for DNSSEC records. Cryptographic
// verification remains a separate trust-layer concern.
type DNSSECRR struct {
	Type       uint16
	Algorithm  uint8
	DigestType uint8
	KeyTag     uint16
	Flags      uint16
	Protocol   uint8
}

// ParseDNSSECMetadata extracts bounded metadata from DS (43), RRSIG (46),
// DNSKEY (48), and NSEC/NSEC3 (47/50) records.
func ParseDNSSECMetadata(msg []byte, typ uint16, start, length int) (DNSSECRR, error) {
	if start < 0 || length < 0 || start+length > len(msg) {
		return DNSSECRR{}, fmt.Errorf("dnssec rdata bounds invalid")
	}
	var out DNSSECRR
	out.Type = typ
	switch typ {
	case 43: // DS: key tag(2), algorithm(1), digest type(1), digest...
		if length < 4 { return out, fmt.Errorf("invalid DS rdata") }
		out.KeyTag = binary.BigEndian.Uint16(msg[start:start+2])
		out.Algorithm = msg[start+2]
		out.DigestType = msg[start+3]
	case 46: // RRSIG: type covered(2), algorithm(1), labels(1), ...
		if length < 4 { return out, fmt.Errorf("invalid RRSIG rdata") }
		out.Algorithm = msg[start+2]
	case 48: // DNSKEY: flags(2), protocol(1), algorithm(1), public key...
		if length < 4 { return out, fmt.Errorf("invalid DNSKEY rdata") }
		out.Flags = binary.BigEndian.Uint16(msg[start:start+2])
		out.Protocol = msg[start+2]
		out.Algorithm = msg[start+3]
	case 47, 50: // NSEC/NSEC3: retain record identity; type bitmap parsing is adapter-specific.
		if length == 0 { return out, fmt.Errorf("empty NSEC rdata") }
	default:
		return out, fmt.Errorf("unsupported DNSSEC RR type: %d", typ)
	}
	return out, nil
}
