package dns

import("context";"fmt";"strings")
type ConsistencyImport struct{Provider ProviderType `json:"provider"`;Zone string `json:"zone"`;Records []Record `json:"records"`}
type ConsistencyImportResult struct{Snapshot ZoneSnapshot `json:"snapshot"`;Report ConsistencyReport `json:"report"`}
type ConsistencyImporter struct{AI *DNSConsistencyAI}
func NewConsistencyImporter(ai *DNSConsistencyAI)*ConsistencyImporter{return &ConsistencyImporter{AI:ai}}
func(i *ConsistencyImporter)Import(ctx context.Context,in ConsistencyImport)(ConsistencyImportResult,error){if i==nil||i.AI==nil{return ConsistencyImportResult{},fmt.Errorf("consistency AI is required")};if ctx==nil{return ConsistencyImportResult{},fmt.Errorf("context is required")};in.Zone=strings.TrimSpace(in.Zone);if in.Zone==""{return ConsistencyImportResult{},fmt.Errorf("zone is required")};if in.Provider==""{return ConsistencyImportResult{},fmt.Errorf("provider is required")};select{case<-ctx.Done():return ConsistencyImportResult{},ctx.Err();default:};snapshot,err:=SnapshotZone(ctx,in.Provider,in.Zone,in.Records);if err!=nil{return ConsistencyImportResult{},err};report,err:=i.AI.Analyze(ctx,map[string][]Record{in.Zone:snapshot.Records});if err!=nil{return ConsistencyImportResult{},err};return ConsistencyImportResult{Snapshot:snapshot,Report:report},nil}
