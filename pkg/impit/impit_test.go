package impit

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateClient(t *testing.T) {
	opts := ImpitOptions{
		Timeout: 5000,
		BaseURL: "https://httpbin.org",
	}
	client := CreateClient(opts)
	if client == nil {
		t.Fatal("Client is nil")
	}
}

func TestHandleRequest(t *testing.T) {
	opts := ImpitOptions{
		BaseURL: "https://httpbin.org",
		Timeout: 10000,
	}
	client := CreateClient(opts)

	reqOpts := RequestInit{
		Method: "GET",
		URL:    "/get",
		QueryParams: map[string]string{
			"foo": "bar",
		},
	}

	resp := HandleRequest(client, reqOpts)

	if resp.Error != "" {
		t.Fatalf("Request failed: %s", resp.Error)
	}

	if resp.Status != 200 {
		t.Errorf("Expected status 200, got %d", resp.Status)
	}

	if !resp.IsBinary {
		t.Error("Expected IsBinary to be true")
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(resp.Body)
	if err != nil {
		t.Fatalf("Failed to decode body: %v", err)
	}

	bodyStr := string(decodedBytes)
	if !strings.Contains(bodyStr, `"foo": "bar"`) && !strings.Contains(bodyStr, `"foo":"bar"`) {
		t.Errorf("Response body does not contain query param. Body: %s", bodyStr)
	}
}

func TestAmazonRequest(t *testing.T) {
	opts := ImpitOptions{
		BaseURL:   "https://www.amazon.com",
		Timeout:   20000,
		Browser:   "chrome",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Headers: map[string]string{
			"accept-language": "en-US,en;q=0.9",
		},
	}
	client := CreateClient(opts)

	reqOpts := RequestInit{
		Method:  "GET",
		URL:     "/",
		Timeout: 20000,
	}

	if reqOpts.Headers == nil {
		reqOpts.Headers = map[string]string{}
	}
	if cookieHeader := os.Getenv("AMAZON_COOKIE"); cookieHeader != "" {
		reqOpts.Headers["Cookie"] = cookieHeader
	}

	resp := HandleRequest(client, reqOpts)
	if resp.Error != "" {
		t.Fatalf("Request failed: %s", resp.Error)
	}

	t.Logf("amazon status=%d statusText=%s finalUrl=%s", resp.Status, resp.StatusText, resp.URL)
	if len(resp.Headers) > 0 {
		if v, ok := resp.Headers["Content-Type"]; ok {
			t.Logf("amazon header Content-Type=%s", v)
		} else if v, ok := resp.Headers["content-type"]; ok {
			t.Logf("amazon header content-type=%s", v)
		}
		if v, ok := resp.Headers["Server"]; ok {
			t.Logf("amazon header Server=%s", v)
		} else if v, ok := resp.Headers["server"]; ok {
			t.Logf("amazon header server=%s", v)
		}
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(resp.Body)
	if err != nil {
		t.Fatalf("Failed to decode body: %v", err)
	}

	bodyStr := string(decodedBytes)
	t.Logf("amazon body_len=%d", len(bodyStr))
	t.Log(bodyStr)

	outPath := filepath.Join("..", "..", "amazon_response.html")
	if absPath, absErr := filepath.Abs(outPath); absErr == nil {
		outPath = absPath
	}
	if writeErr := os.WriteFile(outPath, decodedBytes, 0o644); writeErr != nil {
		t.Fatalf("Failed to write amazon_response.html: %v", writeErr)
	}
	t.Logf("amazon body written to=%s", outPath)
}
