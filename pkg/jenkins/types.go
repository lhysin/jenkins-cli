package jenkins

type Job struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type JobListResponse struct {
	Jobs []Job `json:"jobs"`
}

type Build struct {
	Number    int    `json:"number"`
	URL       string `json:"url"`
	Phase     string `json:"phase,omitempty"`
	Status    string `json:"status,omitempty"`
	Duration  int64  `json:"duration,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Result    string `json:"result,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type BuildListResponse struct {
	Builds []Build `json:"builds"`
}

type Node struct {
	DisplayName string `json:"displayName"`
	Offline     bool   `json:"offline"`
}

type NodeListResponse struct {
	Computer []Node `json:"computer"`
}
