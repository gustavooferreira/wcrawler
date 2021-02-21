package core

// ----------
// Interfaces
// ----------

// Connector describes the connector interface
type Connector interface {
	GetLinks(baseURL string) (statusCode int, links []string, err error)
}
