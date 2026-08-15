package ftnstorage

import "time"

// ScrubSchedule describes bounded integrity checks for stored chunks.
type ScrubSchedule struct {
	Interval time.Duration
	Batch    uint32
}

func (s ScrubSchedule) Valid() bool { return s.Interval > 0 && s.Batch > 0 }
