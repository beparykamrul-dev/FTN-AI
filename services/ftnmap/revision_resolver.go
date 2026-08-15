package ftnmap

// RevisionResolver accepts only newer topology revisions.
type RevisionResolver struct{ Revision uint64 }

func (r *RevisionResolver) Accept(in SyncEnvelope) bool {
	if !in.Valid() || in.Revision <= r.Revision { return false }
	r.Revision = in.Revision
	return true
}
