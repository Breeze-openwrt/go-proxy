# Design: GitHub Action Latest Release

## 1. 编译环境配置 (Build Environment)
- **Golang**: 使用 `actions/setup-go` 指定 1.22+ 版本。
- **静态编译参数**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o sni-proxy main.go`。
- **压缩**: 使用 `upx` (可选) 减小文件体积。

## 2. 单版本发布算法 (Singleton Release Logic)
1. **获取当前 Tag**: 通过系统变量 `${{ github.ref_name }}`。
2. **发布新版本**: 使用 `softprops/action-gh-release`。
3. **清理旧版本 (关键逻辑)**: 
    - 使用 `gh` (GitHub CLI) 执行命令。
    - 列出所有 Release，排除掉刚刚创建的那一个。
    - 循环执行 `gh release delete <old-tag> -y`。
    - 循环执行 `git push --delete origin <old-tag>`。

## 3. 权限要求
- 需要 `GITHUB_TOKEN` 具备 `contents: write` 权限，以执行发布和删除操作。
