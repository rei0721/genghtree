// Package render 实现目录树的渲染输出
package render

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/rei0721/genghtree/internal/domain"
)

const (
	// 树形结构字符
	boxVertical      = "│"    // 垂直线
	boxHorizontal    = "──"   // 水平线
	boxUpAndRight    = "└"    // 拐角
	boxVerticalRight = "├"    // T 形
	boxSpace         = "    " // 空格占位
	boxVerticalSpace = "│   " // 垂直线占位
)

// ASCIIPrinter 实现 domain.Printer 接口的 ASCII 渲染器
type ASCIIPrinter struct {
	writer io.Writer
}

// NewASCIIPrinter 创建一个新的 ASCII 渲染器
func NewASCIIPrinter() *ASCIIPrinter {
	return &ASCIIPrinter{
		writer: os.Stdout,
	}
}

// NewASCIIPrinterWithWriter 创建一个指定输出的 ASCII 渲染器
func NewASCIIPrinterWithWriter(w io.Writer) *ASCIIPrinter {
	return &ASCIIPrinter{
		writer: w,
	}
}

// Print 打印目录树
func (p *ASCIIPrinter) Print(root *domain.Node) error {
	if root == nil {
		return nil
	}

	// 打印根节点名称
	fmt.Fprintln(p.writer, root.Name)

	// 排序子节点：目录在前，文件在后，同类型按名称排序
	sortChildren(root.Children)

	// 递归打印子节点
	p.printChildren(root.Children, "")

	return nil
}

// printChildren 递归打印子节点
func (p *ASCIIPrinter) printChildren(children []*domain.Node, prefix string) {
	for i, child := range children {
		isLast := i == len(children)-1

		// 选择连接符
		connector := boxVerticalRight
		if isLast {
			connector = boxUpAndRight
		}

		// 构建显示名称
		name := child.Name
		if child.IsDir {
			name = name + "/"
		}

		// 打印当前节点
		fmt.Fprintf(p.writer, "%s%s%s %s\n", prefix, connector, boxHorizontal, name)

		// 如果有子节点，递归打印
		if len(child.Children) > 0 {
			// 排序子节点
			sortChildren(child.Children)

			// 计算新的前缀
			newPrefix := prefix
			if isLast {
				newPrefix += boxSpace
			} else {
				newPrefix += boxVerticalSpace
			}

			p.printChildren(child.Children, newPrefix)
		}
	}
}

// sortChildren 对节点进行排序：目录在前，文件在后，同类型按名称排序
func sortChildren(children []*domain.Node) {
	sort.Slice(children, func(i, j int) bool {
		// 目录优先
		if children[i].IsDir != children[j].IsDir {
			return children[i].IsDir
		}
		// 同类型按名称排序
		return children[i].Name < children[j].Name
	})
}

// MarkdownPrinter 实现 Markdown 格式的渲染器
type MarkdownPrinter struct {
	writer     io.Writer
	repoURL    string
	branchName string
}

// NewMarkdownPrinter 创建一个 Markdown 渲染器
func NewMarkdownPrinter(w io.Writer, owner, repo, branch string) *MarkdownPrinter {
	return &MarkdownPrinter{
		writer:     w,
		repoURL:    fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		branchName: branch,
	}
}

// Print 打印 Markdown 格式的目录树
func (p *MarkdownPrinter) Print(root *domain.Node) error {
	if root == nil {
		return nil
	}

	// 打印标题和仓库信息
	fmt.Fprintf(p.writer, "# %s 目录结构\n\n", root.Name)
	fmt.Fprintf(p.writer, "**仓库地址**: [%s](%s)  \n", p.repoURL, p.repoURL)
	fmt.Fprintf(p.writer, "**分支**: `%s`  \n", p.branchName)
	fmt.Fprintf(p.writer, "**生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintln(p.writer, "---")

	// 打印树结构（使用代码块）
	fmt.Fprintln(p.writer, "```")
	fmt.Fprintln(p.writer, root.Name)

	// 排序子节点
	sortChildren(root.Children)

	// 递归打印子节点
	p.printChildren(root.Children, "")

	fmt.Fprintln(p.writer, "```")

	return nil
}

// printChildren 递归打印子节点
func (p *MarkdownPrinter) printChildren(children []*domain.Node, prefix string) {
	for i, child := range children {
		isLast := i == len(children)-1

		// 选择连接符
		connector := boxVerticalRight
		if isLast {
			connector = boxUpAndRight
		}

		// 构建显示名称
		name := child.Name
		if child.IsDir {
			name = name + "/"
		}

		// 打印当前节点
		fmt.Fprintf(p.writer, "%s%s%s %s\n", prefix, connector, boxHorizontal, name)

		// 如果有子节点，递归打印
		if len(child.Children) > 0 {
			// 排序子节点
			sortChildren(child.Children)

			// 计算新的前缀
			newPrefix := prefix
			if isLast {
				newPrefix += boxSpace
			} else {
				newPrefix += boxVerticalSpace
			}

			p.printChildren(child.Children, newPrefix)
		}
	}
}
