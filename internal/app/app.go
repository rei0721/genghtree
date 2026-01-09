// Package app 实现核心业务逻辑
package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/rei0721/genghtree/internal/domain"
)

// App 是应用程序的核心结构
type App struct {
	fetcher domain.Fetcher
	printer domain.Printer
}

// New 创建一个新的 App 实例
func New(fetcher domain.Fetcher, printer domain.Printer) *App {
	return &App{
		fetcher: fetcher,
		printer: printer,
	}
}

// Run 执行主要业务逻辑
func (a *App) Run(ctx context.Context, owner, repo, ref string) error {
	// 1. 获取目录树数据
	treeResp, err := a.fetcher.GetTree(ctx, owner, repo, ref)
	if err != nil {
		return fmt.Errorf("获取仓库目录树失败: %w", err)
	}

	// 2. 检查是否被截断
	if treeResp.Truncated {
		fmt.Println("⚠️  警告: 仓库过大，结果已被截断，显示的目录结构不完整")
	}

	// 3. 构建树结构 (Flat-to-Tree 算法)
	root := a.buildTree(treeResp.Tree, repo)

	// 4. 渲染输出
	return a.printer.Print(root)
}

// buildTree 实现 Flat-to-Tree 算法
// 将扁平的文件路径列表转换为多叉树结构
func (a *App) buildTree(blobs []domain.GitBlob, repoName string) *domain.Node {
	// 创建根节点
	root := domain.NewNode(repoName, "", true)

	// 遍历所有 blob，构建树
	for i := range blobs {
		blob := &blobs[i]
		a.insertPath(root, blob)
	}

	return root
}

// insertPath 将一个路径插入到树中
// 使用前缀树 (Trie) 的思想，逐层查找或创建节点
func (a *App) insertPath(root *domain.Node, blob *domain.GitBlob) {
	parts := strings.Split(blob.Path, "/")
	current := root

	for i, part := range parts {
		if part == "" {
			continue
		}

		// 查找是否已存在该子节点
		child := current.GetChild(part)

		if child == nil {
			// 判断是否为最后一个部分（文件）还是中间部分（目录）
			isDir := blob.Type == "tree" || i < len(parts)-1
			fullPath := strings.Join(parts[:i+1], "/")

			child = domain.NewNode(part, fullPath, isDir)

			// 如果是最后一部分，保存原始元数据
			if i == len(parts)-1 {
				child.Meta = blob
			}

			current.AddChild(child)
		}

		current = child
	}
}
