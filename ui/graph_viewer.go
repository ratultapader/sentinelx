package ui

type GraphNodeDTO struct {
	ID         string                 `json:"id"`
	Label      string                 `json:"label"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type GraphLinkDTO struct {
	Source string                 `json:"source"`
	Target string                 `json:"target"`
	Type   string                 `json:"type"`
	Props  map[string]interface{} `json:"properties,omitempty"`
}

type SecurityGraphDTO struct {
	Nodes []GraphNodeDTO `json:"nodes"`
	Links []GraphLinkDTO `json:"links"`
}