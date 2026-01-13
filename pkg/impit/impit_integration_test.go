package impit

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

// TestHandleRequest_Integration performs a real HTTP request to httpbin.org
// to verify the end-to-end functionality of the impit package.
// This test requires internet access.
func TestHandleRequest_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	opts := ImpitOptions{
		BaseURL: "https://www.amazon.com",
		Timeout: 10000,
		Browser: "chrome",
		Headers: map[string]string{
			"cookie": "session-id=520-3432763-6619032; lc-acbuk=en_GB; ubid-acbuk=261-8790247-0855808; sso-state-acbuk=Xdsso|ZQG-J6z6kqG231syhNdNwhS5_Q7lcebxYKWUYRvz0Svu0x-axk_3b2Rb5NEeGPEAXIjkHc8KUNsUJgDg6MbaRuVXDRZg5ajuU505UwHWg8KZ3Ovr; i18n-prefs=GBP; at-acbuk=Atza|gQC4SOQtAwEBALe-NKeevTEzTSfrCiEdLfTAZiaglHSVfBbK2B_RpLt4rMAMIprXUtKQFv0gCICBSKWRuGpN9a_nAIW_SdgiFXfiREYIUEV19FIHjYgSMpGmgmtdwYqN48lecXt1-71N5FWHUt1o34rRZHkTPyYwfjuaxjHL3r_xGZjP2H5pYX3i5l5w3KP87RAbLPtd_0v1Ao-RkGqHOYEdpXzvBGdreYuwKmHLvznQ1EB6fRvAR6ZJsCBchwqD7o03Qe-ViNuM273Nmxcw2mowfKF5PCwHP1spQjxh9ZGRvnwgPHsQqZwQ7euEw9eOMhWQxg45HvGEFBdT0hrdrqgGQCSNdDFHOKgQT30j3wlobamZ_SFPSEmw1zpfGlZYdhGsk7ymgRSNtspvVkvBrHR9DUDnQZsVFjoVSzPAWoXMDKE; sess-at-acbuk=h8dyKZqGM54ZLiasgvsg9UM7HSVPWcSGBD9l4z7zc2s=; sst-acbuk=Sst1|PQLzwohMykcZpPcYyruIYpxBCivK2t-Y2deQrVmDQpvwd5JAi7BnYGEUZMfxeJm-htVNEJe3GuTmiMFT87mNiHAdTXRRV-dXWq4T2d92y6mzRve0XjDO1oLRZfGweDe9S8lKQB7Lvowwcc2HpXBDnfyvCZKdxcHaMVGrynZUfqaYtZcnUhpRBAHp6sQtEEhOq_KL_ErRQCoAU-Rz2_KhxHnyPxhFDcJw_iPI23ZLvvAdtMcyDAhCVYs5yXdGu_sndL-iYQu5aqdta7lNTzDefH4Wiar69YD0BQe-aC4ybrVOtylJF4hFWNE6xAFH3pYVuJMq; session-id-time=2082787201l; x-acbuk="D6E1EmC2Gt61jieU59NPG@O5RMFD9xKJwxiG4ZsL8om8eOTlzClWtzNoBvbyqP0O"; session-token=GZ9fEfY09afC6U6bj83sUbKH7ODDgAj/dOCPLvuCREFyH8sOapKp5aUSuLwzx3GGyIpr2t8UQQEOTqBLHN1VeUOPx+5JDOJeP9i6rLBzjtfe46Lq/MWgxAnxD4Zo9yqtHV7C9upOC3+2964HqiQYwCdKob6jM4tT7CNmMzB+7ePbhMhBxyDFvfYC4KQowJg2fRr6hkM8vdDy5mZCOxX/T7Gf5z+gAlk/j8O4gnRtot88+AeCdhshuHpiVr2I7XbmCqszH/w9iwjFVXFt6UR6KBAjlTkHaKJYSKxv0ukfR5NfxL7f4zWKhzKxP3LrM4OdPXHpw16xKfL4EEyqHtzu878qvouJdTSEA9eu0I7+h7iEH2e3t97vWE47MFOrzKJC1YTGJPhiSW8=",
		},
		ProxyURL: "http://127.0.0.1:8080",
	}

	client := CreateClient(opts)

	reqInit := RequestInit{
		Method: "GET",
		URL:    "/",
	}

	respData := HandleRequest(client, reqInit)

	if respData.Error != "" {
		t.Fatalf("Request failed with error: %s", respData.Error)
	}

	if respData.Status != 200 {
		t.Errorf("Expected status 200, got %d", respData.Status)
	}

	// Decode body
	bodyBytes, err := base64.StdEncoding.DecodeString(respData.Body)
	if err != nil {
		t.Fatalf("Failed to decode base64 body: %v", err)
	}

	// Parse JSON body
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON body: %v", err)
	}

	// Verify headers reflected by httpbin
	headers, ok := bodyMap["headers"].(map[string]interface{})
	if !ok {
		t.Fatalf("Response JSON does not contain 'headers' object")
	}

	if headers["X-Test-Header"] != "IntegrationTest" {
		t.Errorf("Expected X-Test-Header to be 'IntegrationTest', got %v", headers["X-Test-Header"])
	}

	// Verify User-Agent (Chrome impersonation)
	ua, ok := headers["User-Agent"].(string)
	if !ok || ua == "" {
		t.Errorf("User-Agent missing or empty")
	}
	t.Logf("User-Agent: %s", ua)
}
