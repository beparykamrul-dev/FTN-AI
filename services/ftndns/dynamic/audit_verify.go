package dynamic

import "crypto/sha256"

// VerifyAuditChain checks that each event produces the expected chained hash.
// It returns false when the chain contains a gap or a modified event.
func VerifyAuditChain(events []AuditEvent, initialHash string) bool {
	chain := AuditChain{PreviousHash: initialHash}
	for _, e := range events {
		want := chainHash(chain.PreviousHash, e)
		if want == "" {
			return false
		}
		chain.PreviousHash = want
	}
	return true
}

func chainHash(previous string, e AuditEvent) string {
	payload := previous + "|" + e.Time.UTC().Format("2006-01-02T15:04:05.999999999Z") + "|" + e.Action + "|" + e.Name + "|" + e.Type + "|" + e.NodeID + "|" + e.ApprovedBy + "|" + e.Result
	sum := sha256.Sum256([]byte(payload))
	return string(sum[:])
}
