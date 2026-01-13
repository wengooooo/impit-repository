package impit

// ResponseData mirrors the response structure
type ResponseData struct {
	Status     int               `json:"status"`
	StatusText string            `json:"statusText"`
	Headers    map[string][]string `json:"headers"`
	Body       string            `json:"body"`
	IsBinary   bool              `json:"is_binary"`
	URL        string            `json:"url"`
	Error      string            `json:"error,omitempty"`
}
