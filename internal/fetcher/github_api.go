// Package fetcher 实现 GitHub API 客户端
package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rei0721/genghtree/internal/domain"
)

const (
	// GitHubAPIBase GitHub API 基础 URL
	GitHubAPIBase = "https://api.github.com"
	// DefaultTimeout 默认超时时间
	DefaultTimeout = 30 * time.Second
)

// GitHubClient 实现 domain.Fetcher 接口的 GitHub API 客户端
type GitHubClient struct {
	client  *http.Client
	token   string
	baseURL string
}

// NewGitHubClient 创建一个新的 GitHub API 客户端
func NewGitHubClient(token string) *GitHubClient {
	// 优先使用传入的 token，否则读取环境变量
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	return &GitHubClient{
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
		token:   token,
		baseURL: GitHubAPIBase,
	}
}

// GetTree 获取仓库的目录树
func (c *GitHubClient) GetTree(ctx context.Context, owner, repo, ref string) (*domain.TreeResponse, error) {
	// 构建请求 URL: GET /repos/{owner}/{repo}/git/trees/{ref}?recursive=1
	url := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=1", c.baseURL, owner, repo, ref)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "genghtree-cli")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// 如果有 token，添加认证头
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 GitHub API 失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	var treeResp domain.TreeResponse
	if err := json.Unmarshal(body, &treeResp); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	return &treeResp, nil
}

// handleErrorResponse 处理错误响应
func (c *GitHubClient) handleErrorResponse(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("仓库或分支不存在 (404): %s", string(body))
	case http.StatusForbidden:
		// 检查是否是速率限制
		if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining == "0" {
			resetTime := resp.Header.Get("X-RateLimit-Reset")
			resetUnix, _ := strconv.ParseInt(resetTime, 10, 64)
			resetAt := time.Unix(resetUnix, 0).Format(time.RFC3339)
			return fmt.Errorf("已达到 GitHub API 速率限制，请配置 GITHUB_TOKEN 或等待至 %s 后重试", resetAt)
		}
		return fmt.Errorf("访问被拒绝 (403): %s", string(body))
	case http.StatusUnauthorized:
		return fmt.Errorf("认证失败 (401): 请检查 GITHUB_TOKEN 是否有效")
	default:
		return fmt.Errorf("API 请求失败 (%d): %s", resp.StatusCode, string(body))
	}
}
