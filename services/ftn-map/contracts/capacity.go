package fiber

import "context"

type CoreState struct { Core int; State string; ReservedFor string }
type CableCapacity struct { PathID string; TotalCores int; UsedCores int; SpareCores int; DamagedCores int; ReservedCores int; Cores []CoreState }

type CapacityRepository interface { GetCapacity(context.Context, string) (CableCapacity, error); SaveCapacity(context.Context, CableCapacity) error }
