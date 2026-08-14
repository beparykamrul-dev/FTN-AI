package fiber

import "context"

type OpticalSample struct {
	PathID string
	ONUId string
	TxPowerDbm float64
	RxPowerDbm float64
	AttenuationDb float64
	OpticalBudgetDb float64
	WavelengthNm int
	MeasuredAt string
}

type OpticalRepository interface {
	Observe(context.Context, OpticalSample) error
	Latest(context.Context, string, string) (OpticalSample, error)
}

type OpticalAnalyzer interface {
	Assess(context.Context, OpticalSample) (string, error)
	DetectDegradation(context.Context, []OpticalSample) (bool, error)
}
