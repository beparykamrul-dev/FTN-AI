package ftnstorage

import "time"

// BackupObject represents an immutable recovery point.
type BackupObject struct {
	ID        string
	CreatedAt time.Time
	ExpiresAt time.Time
	Locked    bool
	Verified  bool
}

func (b BackupObject) Usable(now time.Time) bool {
	return b.ID != "" && b.Locked && b.Verified && now.Before(b.ExpiresAt)
}
