package graph

type Elements struct {
	Links []Link `json:"links,omitempty"`
	Nodes []Node `json:"nodes,omitempty"`
}

type Link struct {
	Source string `json:"source,omitempty"`
	Target string `json:"target,omitempty"`
}

type Node struct {
	ID         string `json:"id,omitempty"`
	Domain     string `json:"domain,omitempty"`
	URL        string `json:"url,omitempty"`
	LinksCount int    `json:"linksCount,omitempty"`
}
