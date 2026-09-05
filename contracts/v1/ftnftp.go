package contracts

// FTNFTP is the provider-neutral file/object delivery contract used by FTN.
// It supports classic file-transfer semantics while keeping storage and edge
// implementations replaceable.
type FTNFTPObject struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenant_id"`
	Bucket            string `json:"bucket"`
	Key               string `json:"key"`
	Size              int64  `json:"size"`
	ContentType       string `json:"content_type,omitempty"`
	ETag              string `json:"etag,omitempty"`
	Checksum          string `json:"checksum,omitempty"`
	OriginNode        string `json:"origin_node,omitempty"`
	StorageClass      string `json:"storage_class"`
	ReplicationClass  string `json:"replication_class"`
	Status            string `json:"status"`
}

type FTNFTPTransfer struct {
	ID           string `json:"id"`
	ObjectID     string `json:"object_id"`
	Direction    string `json:"direction"`
	Protocol     string `json:"protocol"`
	EdgeNode     string `json:"edge_node,omitempty"`
	Bytes        int64  `json:"bytes"`
	ThroughputBps int64 `json:"throughput_bps"`
	State        string `json:"state"`
	Error        string `json:"error,omitempty"`
}

type FTNFTPPolicy struct {
	Resume          bool `json:"resume"`
	IntegrityCheck  bool `json:"integrity_check"`
	LocalFirst      bool `json:"local_first"`
	EdgeCache       bool `json:"edge_cache"`
	Replicate       bool `json:"replicate"`
	EncryptAtRest   bool `json:"encrypt_at_rest"`
	Versioning      bool `json:"versioning"`
}
