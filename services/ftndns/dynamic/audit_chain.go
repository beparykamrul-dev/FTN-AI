package dynamic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// AuditChain links FTNDNS audit events so tampering or missing history can be detected.
type AuditChain struct {
	PreviousHash string
}

// Append returns the hash for the next audit entry. The event itself remains
// immutable; persistence and retention are handled by the caller.
func (c *AuditChain) Append(e AuditEvent) string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%d|%s|%s|%s|%s", c.PreviousHash, e.Time.UTC().Format("2006-01-02T15:04:05.999999999Z"), e.Action, e.Name, e.Version, e.NodeID, e.Type, e.ApprovedBy, e.Result)
	sum := sha256.Sum256([]byte(payload))
	h := hex.EncodeToString(sum[:])
	c.PreviousHash = h
	return h
}
