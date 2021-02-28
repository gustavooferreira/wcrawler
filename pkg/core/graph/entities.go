package graph

type Element struct {
	Group string `json:"group,omitempty"`
	Data  Data   `json:"data,omitempty"`
}

type Data struct {
	ID     string `json:"id,omitempty"`
	Parent string `json:"parent,omitempty"`
	Source string `json:"source,omitempty"`
	Target string `json:"target,omitempty"`
}
