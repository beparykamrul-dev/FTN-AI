package ftnstorage

import "time"

// DRDecision describes a disaster-recovery decision without executing it.
type DRDecision struct {
	BackupID string
	Target   string
	Reason   string
	Created  time.Time
	Approved bool
}

func PlanDisasterRecovery(backups []BackupObject, target string, now time.Time) (DRDecision, bool) {
	for i := len(backups) - 1; i >= 0; i-- {
		if backups[i].Usable(now) && target != "" {
			return DRDecision{BackupID: backups[i].ID, Target: target, Reason: "verified-immutable-backup", Created: now.UTC()}, true
		}
	}
	return DRDecision{}, false
}
