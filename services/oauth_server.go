package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// OAuthServer 本地 OAuth 服务器
type OAuthServer struct {
	server    *http.Server
	port      int
	nonce     string
	windowId  string
	isRunning bool
	mu        sync.Mutex
	callback  chan *OAuthResult
}

// OAuthResult OAuth 认证结果
type OAuthResult struct {
	Code     string
	State    string
	Error    string
	Success  bool
}

// NewOAuthServer 创建 OAuth 服务器
func NewOAuthServer() *OAuthServer {
	return &OAuthServer{
		callback: make(chan *OAuthResult, 1),
	}
}

// Start 启动 OAuth 服务器
func (s *OAuthServer) Start() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return "", fmt.Errorf("OAuth server already running")
	}

	// 生成随机端口范围 52779-52879
	port, err := s.findAvailablePort()
	if err != nil {
		return "", fmt.Errorf("find available port failed: %w", err)
	}

	// 生成安全参数
	s.nonce = s.generateNonce()
	s.windowId = s.generateWindowId()
	s.port = port

	// 创建 HTTP 服务器
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.handleCallback)
	mux.HandleFunc("/health", s.handleHealth)

	s.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: mux,
	}

	// 启动服务器
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return "", fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	s.isRunning = true
	go func() {
		defer func() {
			s.mu.Lock()
			s.isRunning = false
			s.mu.Unlock()
		}()

		log.Printf("OAuth server started on 127.0.0.1:%d\n", port)
		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("OAuth server error: %v\n", err)
		}
	}()

	return fmt.Sprintf("http://127.0.0.1:%d/callback", port), nil
}

// Stop 停止 OAuth 服务器
func (s *OAuthServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning || s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("OAuth server shutdown error: %v\n", err)
		return err
	}

	s.isRunning = false
	log.Println("OAuth server stopped")
	return nil
}

// WaitForResult 等待 OAuth 结果
func (s *OAuthServer) WaitForResult(timeout time.Duration) (*OAuthResult, error) {
	select {
	case result := <-s.callback:
		return result, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("OAuth timeout after %v", timeout)
	}
}

// GetNonce 获取 nonce
func (s *OAuthServer) GetNonce() string {
	return s.nonce
}

// GetWindowId 获取 windowId
func (s *OAuthServer) GetWindowId() string {
	return s.windowId
}

// handleCallback 处理 OAuth 回调
func (s *OAuthServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	log.Printf("OAuth callback received: %s %s\n", r.Method, r.URL.String())

	// 只处理 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数
	query := r.URL.Query()
	code := query.Get("code")
	state := query.Get("state")
	errorParam := query.Get("error")

	result := &OAuthResult{
		Code:    code,
		State:   state,
		Error:   errorParam,
		Success: errorParam == "",
	}

	// 发送结果
	select {
	case s.callback <- result:
	default:
		// 避免阻塞
	}

	// 构建重定向 URL
	redirectURL := fmt.Sprintf("codeswitch://auth?nonce=%s&windowId=%s", s.nonce, s.windowId)
	if result.Success && code != "" {
		redirectURL += fmt.Sprintf("&code=%s", code)
	}
	if errorParam != "" {
		redirectURL += fmt.Sprintf("&error=%s", errorParam)
	}

	// 重定向到桌面应用
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// handleHealth 健康检查
func (s *OAuthServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// findAvailablePort 查找可用端口
func (s *OAuthServer) findAvailablePort() (int, error) {
	for port := 52779; port <= 52879; port++ {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port in range 52779-52879")
}

// generateNonce 生成随机 nonce
func (s *OAuthServer) generateNonce() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// 降级方案
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// generateWindowId 生成窗口 ID
func (s *OAuthServer) generateWindowId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()%10000)
}