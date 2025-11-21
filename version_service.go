package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const AppVersion = "v1.0.0"

var UpdateManifestURL = "https://static.bewildcard.com/code-switch/latest.json"

type VersionService struct {
	version       string
	lastCheckTime time.Time
	cachedRelease *ReleaseInfo
}

type R2Manifest struct {
	Version     string            `json:"version"`
	PublishedAt string            `json:"published_at"`
	Body        string            `json:"body"`
	Downloads   map[string]string `json:"downloads"`
}

type ReleaseInfo struct {
	Version     string `json:"version"`
	URL         string `json:"url"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
	HasUpdate   bool   `json:"has_update"`
}

func NewVersionService() *VersionService {
	return &VersionService{version: AppVersion}
}

func (vs *VersionService) CurrentVersion() string {
	return vs.version
}

func (vs *VersionService) CheckForUpdates() (*ReleaseInfo, error) {
	if vs.cachedRelease != nil && time.Since(vs.lastCheckTime) < 5*time.Minute {
		return vs.cachedRelease, nil
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Add cache buster to avoid stale content
	url := fmt.Sprintf("%s?t=%d", UpdateManifestURL, time.Now().Unix())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch update manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update server returned status %d: %s", resp.StatusCode, string(body))
	}

	var manifest R2Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %w", err)
	}

	hasUpdate := compareVersions(vs.version, manifest.Version)
	downloadURL := vs.getDownloadURL(manifest.Downloads)

	releaseInfo := &ReleaseInfo{
		Version:     manifest.Version,
		URL:         downloadURL,
		PublishedAt: manifest.PublishedAt,
		Body:        manifest.Body,
		HasUpdate:   hasUpdate,
	}

	vs.cachedRelease = releaseInfo
	vs.lastCheckTime = time.Now()

	return releaseInfo, nil
}

func (vs *VersionService) getDownloadURL(downloads map[string]string) string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	key := ""
	switch os {
	case "darwin":
		key = "darwin_" + arch
	case "windows":
		key = "windows_amd64" // Currently we only build amd64 for windows
	}

	if url, ok := downloads[key]; ok {
		return url
	}

	// Fallback or empty if not found
	return ""
}

// compareVersions 比较两个版本号，如果 latest > current 返回 true
func compareVersions(current, latest string) bool {
	// 移除 'v' 前缀
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	// 补齐长度
	maxLen := len(currentParts)
	if len(latestParts) > maxLen {
		maxLen = len(latestParts)
	}

	for i := 0; i < maxLen; i++ {
		var currentNum, latestNum int
		if i < len(currentParts) {
			fmt.Sscanf(currentParts[i], "%d", &currentNum)
		}
		if i < len(latestParts) {
			fmt.Sscanf(latestParts[i], "%d", &latestNum)
		}

		if latestNum > currentNum {
			return true
		} else if latestNum < currentNum {
			return false
		}
	}

	return false
}
