package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	// 从环境变量获取Telegram Token
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN environment variable is required")
	}

	// 从环境变量获取API密钥
	apiKeysStr := os.Getenv("API_KEYS")
	if apiKeysStr == "" {
		log.Fatal("API_KEYS environment variable is required")
	}

	// 解析API密钥列表
	apiKeys := strings.Split(apiKeysStr, ",")
	apiKeyMap := make(map[string]bool)
	for _, key := range apiKeys {
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey != "" {
			apiKeyMap[trimmedKey] = true
		}
	}

	if len(apiKeyMap) == 0 {
		log.Fatal("At least one API key is required")
	}

	// 解析目标URL
	target, err := url.Parse("https://api.telegram.org")
	if err != nil {
		log.Fatal(err)
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(target)
	
	// 配置代理错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Proxy error", http.StatusBadGateway)
	}

	// 自定义Director以修改请求
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// 设置正确的Host头
		req.Host = target.Host
	}

	// 健康检查端点
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Telegram Bot API Proxy is running\n")
	})

	// 主要代理处理逻辑
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 验证API Key
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			apiKey = r.URL.Query().Get("api_key")
		}

		if apiKey == "" || !apiKeyMap[apiKey] {
			http.Error(w, "Unauthorized: Invalid or missing API Key", http.StatusUnauthorized)
			return
		}

		// 转发请求到Telegram API
		proxy.ServeHTTP(w, r)
	})

	// 获取端口号，默认为8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Telegram Bot API Proxy server started on port %s", port)
	log.Printf("Health check endpoint: http://localhost:%s/health", port)
	log.Printf("Proxy endpoint: http://localhost:%s/", port)
	log.Printf("API Keys loaded: %d", len(apiKeyMap))
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}