package codec

import (
	"fmt"
	"strings"
)

type Capability struct { ID string `json:"id"`; Class string `json:"class"`; Modes []string `json:"modes,omitempty"`; Hardware bool `json:"hardware,omitempty"`; Lossless bool `json:"lossless,omitempty"`; WorkerIsolated bool `json:"worker_isolated,omitempty"` }
type Job struct { CapabilityID string `json:"capability_id"`; InputURI string `json:"input_uri"`; OutputURI string `json:"output_uri,omitempty"`; Options map[string]string `json:"options,omitempty"`; Preserve bool `json:"preserve_original"` }
type Result struct { JobID string `json:"job_id"`; Status string `json:"status"`; OutputURI string `json:"output_uri,omitempty"`; InputSHA256 string `json:"input_sha256,omitempty"`; OutputSHA256 string `json:"output_sha256,omitempty"`; BytesIn int64 `json:"bytes_in,omitempty"`; BytesOut int64 `json:"bytes_out,omitempty"` }

func (c Capability) Valid() bool { return strings.TrimSpace(c.ID) != "" && strings.TrimSpace(c.Class) != "" }

func (j Job) Valid() error {
	if strings.TrimSpace(j.CapabilityID) == "" { return fmt.Errorf("capability_id is required") }
	if strings.TrimSpace(j.InputURI) == "" { return fmt.Errorf("input_uri is required") }
	if j.OutputURI != "" && strings.TrimSpace(j.OutputURI) == "" { return fmt.Errorf("output_uri is invalid") }
	for k, v := range j.Options {
		if strings.TrimSpace(k) == "" { return fmt.Errorf("option key is invalid") }
		if len(v) > 16384 { return fmt.Errorf("option %q is too large", k) }
	}
	return nil
}

func DefaultCapabilities() []Capability {
	return []Capability{
		{ID:"binary-framing",Class:"transport",Modes:[]string{"stream","datagram"}},
		{ID:"compression",Class:"transfer",Modes:[]string{"lossless"},Lossless:true},
		{ID:"chunking",Class:"transfer",Modes:[]string{"resumable","parallel"},Lossless:true},
		{ID:"deduplication",Class:"transfer",Modes:[]string{"content-addressed"},Lossless:true},
		{ID:"hardware-video-encode",Class:"media",Modes:[]string{"h265"},Hardware:true,WorkerIsolated:true},
		{ID:"media-cut",Class:"media",Modes:[]string{"video-processing"},WorkerIsolated:true},
	}
}
