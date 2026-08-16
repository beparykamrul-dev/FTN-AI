package observability

// BackendCapacity captures runtime ingestion pressure used by the FTN router.
type BackendCapacity struct {
	Name string
	LoadPercent float64
	QueueDepth uint64
	IngestPerSecond float64
	StoragePressurePercent float64
}

func (c BackendCapacity) Score() float64 {
	if c.Name == "" { return 0 }
	s := 100.0 - c.LoadPercent*0.4 - c.StoragePressurePercent*0.4 - float64(c.QueueDepth)/100.0
	if s < 0 { return 0 }
	return s
}
