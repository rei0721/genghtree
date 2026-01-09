// Package domain 定义核心数据模型
package domain

// GitBlob 表示 GitHub API 返回的单个文件/目录信息 (DTO)
type GitBlob struct {
	Path string `json:"path"` // 例如 "cmd/main.go"
	Type string `json:"type"` // "blob" (文件) | "tree" (目录)
	Size int    `json:"size"` // 文件大小 (仅对 blob 有效)
	Sha  string `json:"sha"`  // Git SHA
	Mode string `json:"mode"` // 文件模式
	URL  string `json:"url"`  // API URL
}

// TreeResponse 表示 GitHub Trees API 的响应结构
type TreeResponse struct {
	Sha       string    `json:"sha"`
	URL       string    `json:"url"`
	Tree      []GitBlob `json:"tree"`
	Truncated bool      `json:"truncated"` // 如果为 true，表示结果被截断
}
