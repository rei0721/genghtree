// Package domain 定义核心数据模型
package domain

import "context"

// Node 表示目录树中的一个节点 (多叉树)
type Node struct {
	Name     string           // 节点名称，如 "main.go"
	Path     string           // 完整路径，如 "cmd/main.go"
	IsDir    bool             // 是否为目录
	Children []*Node          // 子节点列表
	Meta     *GitBlob         // 原始元数据引用
	parent   *Node            // 父节点引用 (内部使用)
	childMap map[string]*Node // 子节点映射 (内部使用, 用于快速查找)
}

// NewNode 创建一个新的节点
func NewNode(name, path string, isDir bool) *Node {
	return &Node{
		Name:     name,
		Path:     path,
		IsDir:    isDir,
		Children: make([]*Node, 0),
		childMap: make(map[string]*Node),
	}
}

// AddChild 添加子节点
func (n *Node) AddChild(child *Node) {
	child.parent = n
	n.Children = append(n.Children, child)
	n.childMap[child.Name] = child
}

// GetChild 获取指定名称的子节点
func (n *Node) GetChild(name string) *Node {
	return n.childMap[name]
}

// Fetcher 定义了获取数据的能力接口
// 方便后续 Mock 测试或更换为 GitLab 等其他平台
type Fetcher interface {
	GetTree(ctx context.Context, owner, repo, ref string) (*TreeResponse, error)
}

// Printer 定义了输出的能力接口
// 方便切换为 JSON 输出或彩色输出
type Printer interface {
	Print(root *Node) error
}
