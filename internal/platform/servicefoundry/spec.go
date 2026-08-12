package servicefoundry

import "time"

const SchemaVersion = "ftn-service/v1"

type ServiceSpec struct {
	SchemaVersion string            `json:"schema_version"`
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Version       string            `json:"version"`
	Purpose       string            `json:"purpose"`
	Capabilities  []string          `json:"capabilities,omitempty"`
	Interfaces    []InterfaceSpec   `json:"interfaces,omitempty"`
	Dependencies  []DependencySpec  `json:"dependencies,omitempty"`
	Security      SecuritySpec      `json:"security"`
	Reliability   ReliabilitySpec   `json:"reliability"`
	Resources     ResourceSpec      `json:"resources"`
	Observability ObservabilitySpec `json:"observability"`
	ProviderRules ProviderRules     `json:"provider_rules"`
	CreatedAt     time.Time         `json:"created_at"`
}

type InterfaceSpec struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Target string `json:"target,omitempty"`
}

type DependencySpec struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Required bool   `json:"required"`
	Internal bool   `json:"internal"`
}

type SecuritySpec struct {
	AuthenticationRequired bool     `json:"authentication_required"`
	AuthorizationRequired  bool     `json:"authorization_required"`
	EncryptionRequired     bool     `json:"encryption_required"`
	AuditRequired          bool     `json:"audit_required"`
	RateLimitRequired      bool     `json:"rate_limit_required"`
	SecretIsolation        bool     `json:"secret_isolation"`
	ThreatModelRequired    bool     `json:"threat_model_required"`
	AllowedProtocols       []string `json:"allowed_protocols,omitempty"`
}

type ReliabilitySpec struct {
	MinReplicas       int  `json:"min_replicas"`
	DesiredReplicas   int  `json:"desired_replicas"`
	AutoFailover      bool `json:"auto_failover"`
	AutoRecovery      bool `json:"auto_recovery"`
	DataConsistency   string `json:"data_consistency"`
	GracefulShutdown  bool `json:"graceful_shutdown"`
}

type ResourceSpec struct {
	CPURequestMilli    int64 `json:"cpu_request_milli"`
	MemoryRequestBytes  int64 `json:"memory_request_bytes"`
	StorageBytes        int64 `json:"storage_bytes"`
	MaxConnections      int   `json:"max_connections"`
	MaxMessageBytes     int64 `json:"max_message_bytes"`
}

type ObservabilitySpec struct {
	Metrics bool `json:"metrics"`
	Logs    bool `json:"logs"`
	Traces  bool `json:"traces"`
	Health  bool `json:"health"`
}

type ProviderRules struct {
	ExternalHostedServiceRequired bool   `json:"external_hosted_service_required"`
	PreferredImplementation       string `json:"preferred_implementation"`
	RequiredProtocols             []string `json:"required_protocols,omitempty"`
	AnalysisRequired              bool   `json:"analysis_required"`
	LicenseReviewRequired         bool   `json:"license_review_required"`
	SecurityReviewRequired        bool   `json:"security_review_required"`
}

func NewSpec(id, name, version, purpose string) ServiceSpec {
	return ServiceSpec{
		SchemaVersion: SchemaVersion,
		ID: id,
		Name: name,
		Version: version,
		Purpose: purpose,
		Security: SecuritySpec{
			AuthenticationRequired: true,
			AuthorizationRequired: true,
			EncryptionRequired: true,
			AuditRequired: true,
			RateLimitRequired: true,
			SecretIsolation: true,
			ThreatModelRequired: true,
		},
		Reliability: ReliabilitySpec{
			MinReplicas: 2,
			DesiredReplicas: 2,
			AutoFailover: true,
			AutoRecovery: true,
			GracefulShutdown: true,
		},
		Observability: ObservabilitySpec{Metrics: true, Logs: true, Traces: true, Health: true},
		ProviderRules: ProviderRules{
			ExternalHostedServiceRequired: false,
			PreferredImplementation: "ftn-native",
			AnalysisRequired: true,
			LicenseReviewRequired: true,
			SecurityReviewRequired: true,
		},
	}
}
