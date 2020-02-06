package constants

type ExtractedData []struct {
	StartTime int64    `json:"timestamp"`
	Duration  int64    `json:"duration"`
	BTName    string   `json:"value"`
	TierName  string   `json:"serviceName"`
    NodeName  string   `json:"ipv4"`
    TraceID   string   `json:"traceId"`
    ParentID  string   `json:"parentId,omitempty"`
    SpanID    string   `json:"id"`
}
