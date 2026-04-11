# Design: Harden Release Workflow

## 1. 并发控制
增加 `concurrency` 配置项：
```yaml
concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: true
```

## 2. 清理逻辑升级
使用更严谨的列表获取方式：
```bash
# 获取所有存在的 Release 标签
ALL_RELEASES=$(gh release list --limit 100 --json tagName --jq '.[].tagName')
# 过滤并删除
```

## 3. 编译参数加固
确保静态链接：
```bash
go build -trimpath -ldflags="-s -w -extldflags '-static'" -o sni-proxy main.go
```
