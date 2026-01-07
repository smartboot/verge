package reporter

import (
	"strings"
)

// Reporter handles reporting data to the server
type Reporter struct {
	baseURL string
	token   string
	ready   bool
}

// NewReporter creates a new Reporter instance
func NewReporter(baseURL, token string) *Reporter {
	return &Reporter{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		ready:   true,
	}
}

// SetReady sets the ready status
func (r *Reporter) SetReady(ready bool) {
	r.ready = ready
}
