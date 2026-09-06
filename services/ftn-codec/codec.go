package codec

import (
	"fmt"
	"strings"
)

type Capability struct { ID string `json:"id"`; Class string `json:"class"`; Modes []string `json:"modes,omitempty"`; Hardware bool `json:"hardware,omitempty"`; Lossless bool `json:"lossless,omitempty"`; WorkerIsolated bool `json:"worker_isolated,omitempty"` }
type Job struct { CapabilityID string `json:"capability_id"`; InputURI string `json:"input_uri"`; OutputURI string `json:"output_uri,omitempty"`; Options map[string]string `json:"options,omitempty"`; Preserve bool `json:"preserve_original"` }
type Result struct { JobID string `json:"job_id"`; Status string `json:"status"`; OutputURI string `json:"output_uri,omitempty"`; InputSHA256 string `json:"input_sha256,omitempty"`; OutputSHA256 string `json:"output_sha256,omitempty"`; BytesIn int64 `json:"bytes_in,omitempty"`; BytesOut int64 `json:"bytes_out,omitempty"` }

func (c Capability) Valid() bool {
	id, class := strings.TrimSpace(c.ID), strings.TrimSpace(c.Class)
	if id == "" || class == "" || len(id) > 256 || len(class) > 128 || len(c.Modes) > 32 { return false }
	seen := make(map[string]struct{}, len(c.Modes))
	for _, mode := range c.Modes { mode = strings.TrimSpace(mode); if mode == "" || len(mode) > 128 { return false }; if _, ok := seen[mode]; ok { return false }; seen[mode] = struct{}{} }
	return true
}

func (j Job) Valid() error {
	j.CapabilityID, j.InputURI, j.OutputURI = strings.TrimSpace(j.CapabilityID), strings.TrimSpace(j.InputURI), strings.TrimSpace(j.OutputURI)
	if j.CapabilityID == "" { return fmt.Errorf("capability_id is required") }
	if j.InputURI == "" { return fmt.Errorf("input_uri is required") }
	if len(j.CapabilityID) > 256 || len(j.InputURI) > 4096 || len(j.OutputURI) > 4096 { return fmt.Errorf("job field is too large") }
	if len(j.Options) > 256 { return fmt.Errorf("too many job options") }
	for k, v := range j.Options { k = strings.TrimSpace(k); if k == "" || len(k) > 256 { return fmt.Errorf("option key is invalid") }; if len(v) > 16384 { return fmt.Errorf("option %q is too large", k) } }
	return nil
}

func (r Result) Valid() bool {
	jobID, status, output := strings.TrimSpace(r.JobID), strings.TrimSpace(r.Status), strings.TrimSpace(r.OutputURI)
	return jobID != "" && len(jobID) <= 256 && status != "" && len(status) <= 64 && len(output) <= 4096 && r.BytesIn >= 0 && r.BytesOut >= 0 && r.BytesIn <= 1<<40 && r.BytesOut <= 1<<40
}

func DefaultCapabilities() []Capability { return []Capability{{ID:"binary-framing",Class:"transport",Modes:[]string{"stream","datagram"}},{ID:"compression",Class:"transfer",Modes:[]string{"lossless"},Lossless:true},{ID:"chunking",Class:"transfer",Modes:[]string{"resumable","parallel"},Lossless:true},{ID:"deduplication",Class:"transfer",Modes:[]string{"content-addressed"},Lossless:true},{ID:"hardware-video-encode",Class:"media",Modes:[]string{"h265"},Hardware:true,WorkerIsolated:true},{ID:"media-cut",Class:"media",Modes:[]string{"video-processing"},WorkerIsolated:true}} }
