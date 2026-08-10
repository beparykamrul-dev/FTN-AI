package fiber

type Splitter struct {
	ID string `json:"id"`
	NodeID string `json:"node_id"`
	InputPort string `json:"input_port,omitempty"`
	OutputPorts []string `json:"output_ports,omitempty"`
	Ratio string `json:"ratio,omitempty"`
	LossDB float64 `json:"loss_db,omitempty"`
	Status string `json:"status"`
}

type SplitterConnection struct {
	SplitterID string `json:"splitter_id"`
	Port string `json:"port"`
	TargetNodeID string `json:"target_node_id"`
}

// AttachSplitter adds splitter topology to an existing GIS record without
// changing the underlying persistence implementation.
func AttachSplitter(record *GISRecord, splitter Splitter, connections []SplitterConnection) {
	if record == nil { return }
	for i := range record.Links {
		if record.Links[i].ID == splitter.NodeID { return }
	}
	for _, c := range connections {
		if c.SplitterID != splitter.ID || c.TargetNodeID == "" { continue }
		record.Links = append(record.Links, FiberLink{ID: splitter.ID+":"+c.Port, From: splitter.NodeID, To: c.TargetNodeID, Status: splitter.Status})
	}
}
