package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"impit/pkg/impit"
)

func main() {
	// 1. 配置 Client 选项
	opts := impit.ImpitOptions{
		BaseURL: "https://www.amazon.com/dp/B0DHZRKD5H",
		// Timeout:         10000, // 10s
		Browser: "chrome",
		// RetryCount:      3,
		// RetryWaitTime:   1000,
		RetryOnStatus:   []int{500, 502, 503},
		FollowRedirects: ptr(true),
		Debug:           false,
		ProxyURL:        "http://127.0.0.1:8080",
	}

	// 2. 创建 Client
	client := impit.CreateClient(opts)

	// 3. 准备请求参数
	reqInit := impit.RequestInit{
		Method: "GET",
		URL:    "/",
		Headers: map[string]string{
			"cookie": `session-id=520-3432763-6619032; lc-acbuk=en_GB; ubid-acbuk=261-8790247-0855808; sso-state-acbuk=Xdsso|ZQG-J6z6kqG231syhNdNwhS5_Q7lcebxYKWUYRvz0Svu0x-axk_3b2Rb5NEeGPEAXIjkHc8KUNsUJgDg6MbaRuVXDRZg5ajuU505UwHWg8KZ3Ovr; i18n-prefs=GBP; at-acbuk=Atza|gQC4SOQtAwEBALe-NKeevTEzTSfrCiEdLfTAZiaglHSVfBbK2B_RpLt4rMAMIprXUtKQFv0gCICBSKWRuGpN9a_nAIW_SdgiFXfiREYIUEV19FIHjYgSMpGmgmtdwYqN48lecXt1-71N5FWHUt1o34rRZHkTPyYwfjuaxjHL3r_xGZjP2H5pYX3i5l5w3KP87RAbLPtd_0v1Ao-RkGqHOYEdpXzvBGdreYuwKmHLvznQ1EB6fRvAR6ZJsCBchwqD7o03Qe-ViNuM273Nmxcw2mowfKF5PCwHP1spQjxh9ZGRvnwgPHsQqZwQ7euEw9eOMhWQxg45HvGEFBdT0hrdrqgGQCSNdDFHOKgQT30j3wlobamZ_SFPSEmw1zpfGlZYdhGsk7ymgRSNtspvVkvBrHR9DUDnQZsVFjoVSzPAWoXMDKE; sess-at-acbuk=h8dyKZqGM54ZLiasgvsg9UM7HSVPWcSGBD9l4z7zc2s=; sst-acbuk=Sst1|PQLzwohMykcZpPcYyruIYpxBCivK2t-Y2deQrVmDQpvwd5JAi7BnYGEUZMfxeJm-htVNEJe3GuTmiMFT87mNiHAdTXRRV-dXWq4T2d92y6mzRve0XjDO1oLRZfGweDe9S8lKQB7Lvowwcc2HpXBDnfyvCZKdxcHaMVGrynZUfqaYtZcnUhpRBAHp6sQtEEhOq_KL_ErRQCoAU-Rz2_KhxHnyPxhFDcJw_iPI23ZLvvAdtMcyDAhCVYs5yXdGu_sndL-iYQu5aqdta7lNTzDefH4Wiar69YD0BQe-aC4ybrVOtylJF4hFWNE6xAFH3pYVuJMq; session-id-time=2082787201l; x-acbuk="D6E1EmC2Gt61jieU59NPG@O5RMFD9xKJwxiG4ZsL8om8eOTlzClWtzNoBvbyqP0O"; session-token=GZ9fEfY09afC6U6bj83sUbKH7ODDgAj/dOCPLvuCREFyH8sOapKp5aUSuLwzx3GGyIpr2t8UQQEOTqBLHN1VeUOPx+5JDOJeP9i6rLBzjtfe46Lq/MWgxAnxD4Zo9yqtHV7C9upOC3+2964HqiQYwCdKob6jM4tT7CNmMzB+7ePbhMhBxyDFvfYC4KQowJg2fRr6hkM8vdDy5mZCOxX/T7Gf5z+gAlk/j8O4gnRtot88+AeCdhshuHpiVr2I7XbmCqszH/w9iwjFVXFt6UR6KBAjlTkHaKJYSKxv0ukfR5NfxL7f4zWKhzKxP3LrM4OdPXHpw16xKfL4EEyqHtzu878qvouJdTSEA9eu0I7+h7iEH2e3t97vWE47MFOrzKJC1YTGJPhiSW8="`,
		},
	}

	// 4. 发起请求
	respData := impit.HandleRequest(client, reqInit)

	// 5. 处理结果
	handleResponse(respData)

}

func handleResponse(respData impit.ResponseData) {
	if respData.Error != "" {
		log.Printf("Request failed: %s", respData.Error)
		return
	}

	fmt.Printf("Status: %d %s\n", respData.Status, respData.StatusText)
	fmt.Printf("Content-Type: %s\n", respData.Headers["Content-Type"])

	// 6. 解码 Body (impit 默认返回 Base64 编码的 Body 字符串)
	bodyBytes, err := base64.StdEncoding.DecodeString(respData.Body)
	if err != nil {
		fmt.Printf("Body (Base64): %s\n", respData.Body)
		return
	}

	// 7. 尝试解析为 JSON，失败则作为普通字符串打印
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
		headers := bodyMap["headers"]
		fmt.Printf("Response Headers reflected by server: %+v\n", headers)
	} else {
		// 非 JSON 内容（如 HTML），直接打印字符串（截断以防过长）
		content := string(bodyBytes)
		if len(content) > 200 {
			fmt.Printf("Body (String, truncated): %s...\n", content[:200])
		} else {
			fmt.Printf("Body (String): %s\n", content)
		}
	}
}

func ptr(b bool) *bool {
	return &b
}
