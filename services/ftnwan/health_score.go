package ftnwan

// HealthScore produces a bounded 0..100 score from normalized path metrics.
func HealthScore(latencyMS, lossPPM uint32, hops uint16) uint8 {
	score := 100
	if latencyMS > 20 { score -= int((latencyMS - 20) / 2) }
	if lossPPM > 100 { score -= int((lossPPM - 100) / 20) }
	if hops > 1 { score -= int(hops - 1) * 3 }
	if score < 0 { return 0 }
	if score > 100 { return 100 }
	return uint8(score)
}
