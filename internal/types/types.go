package types

type IOC struct {
	ID              string   `json:"id"`
	ThreatType      string   `json:"threat_type"`
	MalwareFamily   string   `json:"malware_family"`
	IOC             string   `json:"ioc"`
	IOCType         string   `json:"ioc_type"`
	ConfidenceLevel int      `json:"confidence_level"`
	Reference       string   `json:"reference"`
	Tags            []string `json:"tags"`
	Timestamp       string   `json:"timestamp"`
}
