package mesh

import "testing"

func TestTopologyRejectsNilReceiver(t *testing.T) {
	var topo *Topology
	if err := topo.UpsertNode(Node{ID:"n1"}); err == nil { t.Fatal("nil topology must fail closed") }
}

func TestNewTopologyInitializesRegistry(t *testing.T) {
	if NewTopology() == nil { t.Fatal("constructor returned nil") }
}
