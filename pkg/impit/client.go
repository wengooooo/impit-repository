package impit

import (
	"time"

	"github.com/imroc/req/v3"
)

// ImpitOptions mirrors the configuration options from impit-node
type ImpitOptions struct {
	Browser            string            `json:"browser,omitempty"` // "chrome" | "firefox"
	IgnoreTLSErrors    bool              `json:"ignore_tls_errors,omitempty"`
	ProxyURL           string            `json:"proxy_url,omitempty"`
	Timeout            int               `json:"timeout,omitempty"` // in milliseconds
	Http3              bool              `json:"http3,omitempty"`
	FollowRedirects    *bool             `json:"follow_redirects,omitempty"` // Use pointer to distinguish default (true)
	MaxRedirects       int               `json:"max_redirects,omitempty"`
	Headers            map[string]string `json:"headers,omitempty"`
	Debug              bool              `json:"debug,omitempty"`
	BaseURL            string            `json:"base_url,omitempty"`
	UserAgent          string            `json:"user_agent,omitempty"`
	ForceHTTP1         bool              `json:"force_http1,omitempty"`
	ForceHTTP2         bool              `json:"force_http2,omitempty"`
	DisableKeepAlives  bool              `json:"disable_keep_alives,omitempty"`
	DisableCompression bool              `json:"disable_compression,omitempty"`
	RetryCount         int               `json:"retry_count,omitempty"`
	RetryWaitTime      int               `json:"retry_wait_time,omitempty"` // milliseconds
	RetryOnStatus      []int             `json:"retry_on_status,omitempty"`
	DisableCookieJar   bool              `json:"disable_cookie_jar,omitempty"`
}

func CreateClient(opts ImpitOptions) *req.Client {
	c := req.C()

	if opts.DisableCookieJar {
		c.SetCookieJar(nil)
	}

	if opts.BaseURL != "" {
		c.SetBaseURL(opts.BaseURL)
	}

	if opts.UserAgent != "" {
		c.SetUserAgent(opts.UserAgent)
	}

	if opts.ForceHTTP1 {
		c.EnableForceHTTP1()
	}

	if opts.ForceHTTP2 {
		c.EnableForceHTTP2()
	}

	if opts.DisableKeepAlives {
		c.DisableKeepAlives()
	}

	if opts.DisableCompression {
		c.DisableCompression()
	}

	if opts.Browser == "chrome" {
		c.ImpersonateChrome()
	} else {
		// Default to Firefox for backward compatibility if empty or unknown
		c.ImpersonateFirefox()
	}

	if opts.Http3 {
		c.EnableHTTP3()
	}

	if opts.IgnoreTLSErrors {
		c.EnableInsecureSkipVerify()
	}

	if opts.ProxyURL != "" {
		c.SetProxyURL(opts.ProxyURL)
	}

	if opts.Timeout > 0 {
		c.SetTimeout(time.Duration(opts.Timeout) * time.Millisecond)
	} else {
		c.SetTimeout(30 * time.Second)
	}

	if opts.Debug {
		c.EnableDumpAll()
	}

	// req/v3 follows redirects by default. Only disable if explicitly set to false.
	if opts.FollowRedirects != nil && !*opts.FollowRedirects {
		c.SetRedirectPolicy(req.NoRedirectPolicy())
	} else if opts.MaxRedirects > 0 {
		c.SetRedirectPolicy(req.MaxRedirectPolicy(opts.MaxRedirects))
	}

	if len(opts.Headers) > 0 {
		c.SetCommonHeaders(opts.Headers)
	}

	if opts.RetryCount > 0 {
		c.SetCommonRetryCount(opts.RetryCount)
		if opts.RetryWaitTime > 0 {
			c.SetCommonRetryFixedInterval(time.Duration(opts.RetryWaitTime) * time.Millisecond)
		}

		// Configure retry condition if specific status codes are provided
		if len(opts.RetryOnStatus) > 0 {
			c.SetCommonRetryCondition(func(resp *req.Response, err error) bool {
				// Always retry on network errors (default behavior)
				if err != nil {
					return true
				}
				// Retry on specified status codes
				if resp != nil {
					for _, status := range opts.RetryOnStatus {
						if resp.StatusCode == status {
							return true
						}
					}
				}
				return false
			})
		}

		// Add retry hook for debugging
		if opts.Debug {
			c.SetCommonRetryHook(func(resp *req.Response, err error) {
				if err != nil {
					c.GetLogger().Warnf("Retrying due to error: %v", err)
				} else if resp != nil {
					c.GetLogger().Warnf("Retrying due to status: %d", resp.StatusCode)
				}
			})
		}
	}

	return c
}
