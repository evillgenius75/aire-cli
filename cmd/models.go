package main

// Model represents a model.
type Model struct {
	Name string `json:"name"`
}

// ModelServer represents a model server.
type ModelServer struct {
	Name string `json:"name"`
}

// ModelServerVersion represents a version of a model server.
type ModelServerVersion struct {
	Name string `json:"name"`
}

// Accelerator represents an accelerator.
type Accelerator struct {
	Name string `json:"name"`
}

// Manifest represents an optimized manifest.
type Manifest struct {
	AcceleratorType string `json:"acceleratorType"`
}

// ModelsAndServers represents a combination of a model and a model server.
type ModelsAndServers struct {
	ModelName       string `json:"modelName"`
	ModelServerName string `json:"modelServerName"`
	CreateTime      string `json:"createTime"`
	UpdateTime      string `json:"updateTime"`
}

// ListModelsResponse represents the response from the ListModels API.
type ListModelsResponse struct {
	ModelNames []string `json:"modelNames"`
}

// ListModelServersResponse represents the response from the ListModelServers API.
type ListModelServersResponse struct {
	ModelServerNames []string `json:"modelServerNames"`
}

// ListModelServerVersionsResponse represents the response from the ListModelServerVersions API.
type ListModelServerVersionsResponse struct {
	ModelServerVersions []string `json:"modelServerVersions"`
}

// ModelAndModelServerInfo represents the modelAndModelServerInfo object.
type ModelAndModelServerInfo struct {
	ModelName          string `json:"modelName"`
	ModelServerName    string `json:"modelServerName"`
	ModelServerVersion string `json:"modelServerVersion"`
}

// ResourcesUsed represents the resourcesUsed object.
type ResourcesUsed struct {
	AcceleratorCount int `json:"acceleratorCount"`
}

// PerformanceStats represents the performanceStats object.
type PerformanceStats struct {
	TpotMilliseconds        int `json:"tpotMilliseconds"`
	QueriesPerSecond        int `json:"queriesPerSecond"`
	OutputTokensPerSecond   int `json:"outputTokensPerSecond"`
	NtpotMilliseconds       int `json:"ntpotMilliseconds"`
}

// AcceleratorOption represents an object in the acceleratorOptions array.
type AcceleratorOption struct {
	AcceleratorType         string                  `json:"acceleratorType"`
	ModelAndModelServerInfo ModelAndModelServerInfo `json:"modelAndModelServerInfo"`
	MachineType             string                  `json:"machineType,omitempty"`
	TpuTopology             string                  `json:"tpuTopology,omitempty"`
	ResourcesUsed           ResourcesUsed           `json:"resourcesUsed"`
	PerformanceStats        PerformanceStats        `json:"performanceStats"`
}

// ListAcceleratorsResponse represents the response from the ListAccelerators API.
type ListAcceleratorsResponse struct {
	AcceleratorOptions          []AcceleratorOption `json:"acceleratorOptions"`
	MinTpotMilliseconds         int                 `json:"minTpotMilliseconds"`
	MaxTpotMilliseconds         int                 `json:"maxTpotMilliseconds"`
	MinThroughputTokensPerSecond int                 `json:"minThroughputTokensPerSecond"`
	MaxThroughputTokensPerSecond int                 `json:"maxThroughputTokensPerSecond"`
	MinNtpotMilliseconds        int                 `json:"minNtpotMilliseconds"`
	MaxNtpotMilliseconds        int                 `json:"maxNtpotMilliseconds"`
}

// K8sManifest represents a Kubernetes manifest.
type K8sManifest struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Content    string `json:"content"`
}

// CreateManifestResponse represents the response from the CreateManifest API.
type CreateManifestResponse struct {
	K8sManifests []K8sManifest `json:"k8sManifests"`
	Comments     []string      `json:"comments"`
}
