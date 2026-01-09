# GengHTree - GitHub Remote Tree Viewer

一个快速查看 GitHub 仓库目录结构的 CLI 工具，无需 clone 仓库。

## 功能特点

- ✅ **零依赖本地仓库** - 无需 clone 代码
- ✅ **快速响应** - 通过 API 直接获取，毫秒级响应
- ✅ **ASCII 树状图** - 终端可视化目录结构
- ✅ **分支/Tag 支持** - 可指定任意分支或标签
- ✅ **Token 认证** - 支持私有仓库和更高的 API 配额
- ✅ **Markdown 输出** - 可保存为 Markdown 文件，包含仓库信息

## 安装

```bash
go install github.com/rei0721/genghtree/cmd/genghtree@latest
```

或直接编译：

```bash
git clone https://github.com/rei0721/genghtree.git
cd genghtree
go build -o genghtree ./cmd/ghtree
```

## 使用方法

```bash
# 查看公开仓库
genghtree owner/repo

# 指定分支或 Tag
genghtree owner/repo --branch develop
genghtree owner/repo -b v1.0.0

# 保存为 Markdown 文件
genghtree owner/repo -o tree.md
genghtree owner/repo -b master -o output.md

# 使用 Token 访问私有仓库或提高 API 配额
genghtree owner/repo --token ghp_xxxx
# 或设置环境变量
export GITHUB_TOKEN=ghp_xxxx
genghtree owner/repo
```

## 示例输出

### 终端输出

```
cobra
├── .github/
│   └── workflows/
│       └── test.yml
├── args.go
├── bash_completions.go
├── cobra.go
├── command.go
└── README.md
```

### Markdown 文件输出

```markdown
# repo-name 目录结构

**仓库地址**: [https://github.com/owner/repo](https://github.com/owner/repo)  
**分支**: `main`  
**生成时间**: 2026-01-09 14:20:25

---

(目录树内容)
```

## 配置

### 环境变量

| 变量名         | 说明                                                      |
| -------------- | --------------------------------------------------------- |
| `GITHUB_TOKEN` | GitHub Personal Access Token，用于私有仓库或提高 API 限额 |

### 命令行参数

| 参数       | 简写 | 默认值 | 说明                 |
| ---------- | ---- | ------ | -------------------- |
| `--branch` | `-b` | `main` | 指定分支或 Tag       |
| `--token`  | `-t` | -      | GitHub Token         |
| `--output` | `-o` | -      | 输出到 Markdown 文件 |

## 限制

- GitHub 匿名 API 限制为 60 次/小时，使用 Token 可提升至 5000 次/小时
- 大型仓库 (如 Linux Kernel) 可能会被 API 截断，程序会显示警告

## License

MIT
