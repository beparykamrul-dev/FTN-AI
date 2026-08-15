package ftnstorage

import "time"

// RecoveryBackoff prevents a repeatedly failing node from entering a tight
// repair loop and gives the control plane bounded retry timing.
type RecoveryBackoff struct {
	Base time.Duration
	Max  time.Duration
}

func (b RecoveryBackoff) Delay(failures uint32) time.Duration {
	if failures == 0 || b.Base <= 0 {
		return 0
	}
	d := b.Base
	for i := uint32(1); i < failures && d < b.Max; i++ {
		d *= 2
	}
	if b.Max > 0 && d > b.Max {
		return b.Max
	}
	return d
}
