// genghtree 是一个 GitHub 远程仓库目录树查看工具
// 无需 clone 仓库即可快速查看目录结构
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rei0721/genghtree/internal/app"
	"github.com/rei0721/genghtree/internal/fetcher"
	"github.com/rei0721/genghtree/internal/render"
)

var (
	// 命令行参数
	branch string
	token  string
	output string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "genghtree <owner/repo>",
		Short: "查看 GitHub 仓库的目录结构",
		Long: `GHTree - GitHub Remote Tree Viewer

一个快速查看 GitHub 仓库目录结构的 CLI 工具。
无需 clone 仓库，直接通过 API 获取并展示目录树。

示例:
  genghtree rei0721/tools
  genghtree rei0721/tools --branch develop
  genghtree rei0721/tools -b v1.0.0
  genghtree rei0721/tools -o tree.md`,
		Args: cobra.ExactArgs(1),
		RunE: runTree,
	}

	// 添加命令行参数
	rootCmd.Flags().StringVarP(&branch, "branch", "b", "main", "指定分支或 Tag")
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "GitHub Personal Access Token (也可设置 GITHUB_TOKEN 环境变量)")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "输出到 Markdown 文件 (例如: tree.md)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runTree(cmd *cobra.Command, args []string) error {
	// 解析 owner/repo
	repoArg := args[0]
	parts := strings.Split(repoArg, "/")
	if len(parts) != 2 {
		return fmt.Errorf("无效的仓库格式，请使用 owner/repo 格式，例如：rei0721/tools")
	}

	owner := parts[0]
	repo := parts[1]

	// 初始化 GitHub 客户端
	githubClient := fetcher.NewGitHubClient(token)

	// 根据是否指定输出文件选择渲染器
	if output != "" {
		// 输出到 Markdown 文件
		file, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("无法创建输出文件: %w", err)
		}
		defer file.Close()

		mdPrinter := render.NewMarkdownPrinter(file, owner, repo, branch)
		application := app.New(githubClient, mdPrinter)

		ctx := context.Background()
		if err := application.Run(ctx, owner, repo, branch); err != nil {
			return err
		}

		fmt.Printf("✅ 目录树已保存到 %s\n", output)
	} else {
		// 输出到终端
		asciiPrinter := render.NewASCIIPrinter()
		application := app.New(githubClient, asciiPrinter)

		ctx := context.Background()
		if err := application.Run(ctx, owner, repo, branch); err != nil {
			return err
		}
	}

	return nil
}
