package impit

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/imroc/req/v3"
)

func HandleRequest(client *req.Client, opts RequestInit) ResponseData {

	fmt.Println("[Go] HandleRequest called")
	fmt.Printf("[Go] Requesting URL: %s\n", opts.URL)
	r := client.R()

	if opts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.Timeout)*time.Millisecond)
		defer cancel()
		r.SetContext(ctx)
	}

	if len(opts.Headers) > 0 {
		r.SetHeaders(opts.Headers)
	}

	if opts.UserAgent != "" {
		r.SetHeader("User-Agent", opts.UserAgent)
	}

	if len(opts.QueryParams) > 0 {
		r.SetQueryParams(opts.QueryParams)
	}

	if len(opts.PathParams) > 0 {
		r.SetPathParams(opts.PathParams)
	}

	if len(opts.FormData) > 0 {
		r.SetFormData(opts.FormData)
	}

	if len(opts.Cookies) > 0 {
		var cookies []*http.Cookie
		for k, v := range opts.Cookies {
			cookies = append(cookies, &http.Cookie{Name: k, Value: v})
		}
		r.SetCookies(cookies...)
	}

	if opts.BasicAuthUser != "" || opts.BasicAuthPass != "" {
		r.SetBasicAuth(opts.BasicAuthUser, opts.BasicAuthPass)
	}

	if opts.BearerToken != "" {
		r.SetBearerAuthToken(opts.BearerToken)
	}

	if opts.BodyBase64 != "" {
		data, err := base64.StdEncoding.DecodeString(opts.BodyBase64)
		if err == nil {
			r.SetBodyBytes(data)
		}
	} else if opts.Body != "" {
		r.SetBody(opts.Body)
	}

	method := "GET"
	if opts.Method != "" {
		method = opts.Method
	}

	resp, err := r.Send(method, opts.URL)
	if err != nil {
		return ResponseData{Error: err.Error()}
	}

	// Read Body
	bodyBytes, err := resp.ToBytes()
	if err != nil {
		return ResponseData{Error: "Failed to read body: " + err.Error()}
	}

	// Convert headers
	headers := make(map[string][]string)
	for k, v := range resp.Header {
		headers[k] = v
	}

	// Encode body to base64 to ensure binary safety
	bodyBase64 := base64.StdEncoding.EncodeToString(bodyBytes)

	return ResponseData{
		Status:     resp.StatusCode,
		StatusText: resp.Status,
		Headers:    headers,
		Body:       bodyBase64,
		IsBinary:   true,
		URL:        resp.Request.URL.String(),
	}
}
