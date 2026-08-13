package adapters

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// ParseExtendedRData decodes the common domain-name based RDATA types that
// the minimal RR parser intentionally leaves opaque: MX, SOA and CAA.
func ParseExtendedRData(msg []byte, typ uint16, start, length int) (string, error) {
	if start < 0 || length < 0 || start+length > len(msg) {
		return "", fmt.Errorf("rdata bounds invalid")
	}
	switch typ {
	case 15: // MX
		if length < 3 { return "", fmt.Errorf("invalid MX rdata") }
		pref := binary.BigEndian.Uint16(msg[start : start+2])
		host, _, err := readName(msg, start+2)
		if err != nil { return "", err }
		return fmt.Sprintf("%d %s", pref, host), nil
	case 6: // SOA
		mname, off, err := readName(msg, start)
		if err != nil { return "", err }
		rname, off, err := readName(msg, off)
		if err != nil || off+20 > start+length { return "", fmt.Errorf("invalid SOA rdata") }
		serial := binary.BigEndian.Uint32(msg[off : off+4])
		refresh := binary.BigEndian.Uint32(msg[off+4 : off+8])
		retry := binary.BigEndian.Uint32(msg[off+8 : off+12])
		expire := binary.BigEndian.Uint32(msg[off+12 : off+16])
		minimum := binary.BigEndian.Uint32(msg[off+16 : off+20])
		return fmt.Sprintf("%s %s %d %d %d %d %d", mname, rname, serial, refresh, retry, expire, minimum), nil
	case 257: // CAA
		if length < 2 { return "", fmt.Errorf("invalid CAA rdata") }
		flags := msg[start]
		tagLen := int(msg[start+1])
		if 2+tagLen > length { return "", fmt.Errorf("invalid CAA tag") }
		tag := string(msg[start+2 : start+2+tagLen])
		value := string(msg[start+2+tagLen : start+length])
		return fmt.Sprintf("%d %s %s", flags, strings.TrimSpace(tag), value), nil
	default:
		return "", fmt.Errorf("unsupported extended RR type: %d", typ)
	}
}
