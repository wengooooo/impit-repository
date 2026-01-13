package impit

// RequestInit mirrors the request options from impit-node
type RequestInit struct {
	Method        string            `json:"method,omitempty"`
	URL           string            `json:"url"`
	Headers       map[string]string `json:"headers,omitempty"`
	UserAgent     string            `json:"user_agent,omitempty"`
	Body          string            `json:"body,omitempty"`
	BodyBase64    string            `json:"body_base64,omitempty"`
	Timeout       int               `json:"timeout,omitempty"`
	ForceHttp3    bool              `json:"force_http3,omitempty"`
	QueryParams   map[string]string `json:"query_params,omitempty"`
	PathParams    map[string]string `json:"path_params,omitempty"`
	Cookies       map[string]string `json:"cookies,omitempty"`
	FormData      map[string]string `json:"form_data,omitempty"`
	BasicAuthUser string            `json:"basic_auth_user,omitempty"`
	BasicAuthPass string            `json:"basic_auth_pass,omitempty"`
	BearerToken   string            `json:"bearer_token,omitempty"`
}
