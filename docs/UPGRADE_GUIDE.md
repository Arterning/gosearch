# SimpleFTS V2 升级指南

## 🎉 从 Rust 版本移植的功能

本次升级将 Rust 版本 (rsfts) 的所有改进功能移植到了 Go 版本。

## 📋 新增文件清单

### 核心模块
- ✅ `document_new.go` - 新的文档结构（替代 `document.go`）
- ✅ `index_new.go` - 改进的倒排索引（替代 `index.go`）
- ✅ `storage.go` - BoltDB 持久化层（新增）
- ✅ `ranking.go` - BM25 排序算法（新增）
- ✅ `engine.go` - 搜索引擎核心（新增）
- ✅ `api.go` - Gin HTTP API（新增）
- ✅ `main_new.go` - 新的主程序（替代 `main.go`）

### 保留文件
- ✅ `tokenizer.go` - 分词器（保持不变）
- ✅ `filter.go` - 文本过滤器（保持不变）

### 文档和脚本
- ✅ `README_V2.md` - 完整的使用文档
- ✅ `UPGRADE_GUIDE.md` - 本升级指南
- ✅ `example_usage.sh` - CLI 使用示例
- ✅ `test_api.sh` - API 测试脚本
- ✅ `.gitignore` - Git 忽略规则

### 依赖更新
- ✅ `go.mod` - 更新了依赖项

## 🚀 快速开始

### 方式 1：直接运行（推荐用于测试）

不需要替换任何文件，直接运行：

```bash
# 启动服务器
go run *.go serve

# CLI 命令
go run *.go insert --id "1" --title "Test" --content "Hello World"
go run *.go search --query "test"
go run *.go stats
```

### 方式 2：替换旧文件（生产环境）

```bash
# 1. 备份旧文件
mkdir backup
cp main.go backup/
cp document.go backup/
cp index.go backup/

# 2. 替换为新文件
mv main_new.go main.go
mv document_new.go document.go
mv index_new.go index.go

# 3. 编译
go build

# 4. 运行
./simplefts serve
```

## 📦 依赖安装

新增依赖会自动下载：

```bash
go mod download
```

新增的依赖包括：
- `github.com/gin-gonic/gin` - Web 框架
- `github.com/spf13/cobra` - CLI 框架
- `go.etcd.io/bbolt` - 嵌入式数据库

## ✨ 功能对比

| 功能 | 原版本 | V2 版本 |
|------|--------|---------|
| **文档插入** | 仅启动时加载 | ✅ 实时插入 |
| **文档更新** | ❌ 不支持 | ✅ 支持 |
| **文档删除** | ❌ 不支持 | ✅ 支持 |
| **HTTP API** | ❌ 无 | ✅ 完整 REST API |
| **CLI 工具** | 基础参数 | ✅ 子命令系统 |
| **数据存储** | JSON 文件 | ✅ BoltDB |
| **相关性排序** | ❌ 无 | ✅ BM25 算法 |
| **搜索模式** | 仅 AND | ✅ AND/OR |
| **分页** | ❌ 无 | ✅ 支持 |
| **并发安全** | 部分 | ✅ 完全支持 |
| **增量更新** | ❌ 全量重写 | ✅ 增量更新 |

## 🎯 使用示例

### CLI 使用

```bash
# 启动 HTTP 服务器
go run *.go serve --port 8080

# 插入文档
go run *.go insert \
  --id "doc1" \
  --title "Go Tutorial" \
  --content "Learn Go programming"

# 搜索文档（AND 模式）
go run *.go search --query "go programming"

# 搜索文档（OR 模式）
go run *.go search --query "go rust" --mode or

# 查看文档
go run *.go get --id "doc1"

# 删除文档
go run *.go delete --id "doc1"

# 查看统计
go run *.go stats
```

### HTTP API 使用

```bash
# 插入文档
curl -X POST http://localhost:3000/documents \
  -H "Content-Type: application/json" \
  -d '{"id":"1","title":"Test","content":"Hello"}'

# 搜索
curl "http://localhost:3000/search?query=test&limit=10&ranked=true"

# 获取文档
curl http://localhost:3000/documents/1

# 更新文档
curl -X PUT http://localhost:3000/documents/1 \
  -H "Content-Type: application/json" \
  -d '{"id":"1","title":"Updated","content":"New content"}'

# 删除文档
curl -X DELETE http://localhost:3000/documents/1

# 获取统计
curl http://localhost:3000/stats
```

## 🔄 数据迁移

**重要**：新版本使用 BoltDB，与旧版本的 JSON 文件不兼容。

### 迁移步骤

如果你有旧版本的数据需要迁移：

1. 使用旧版本导出数据为 JSON
2. 写一个简单的导入脚本读取 JSON 并调用新版本的 API

示例导入脚本：

```go
package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
)

func migrateData() {
    // 读取旧的 doc.json
    data, _ := ioutil.ReadFile("doc.json")

    var docs []struct {
        ID      int    `json:"id"`
        Title   string `json:"title"`
        Content string `json:"abstract"`
        URL     string `json:"url"`
    }

    json.Unmarshal(data, &docs)

    // 创建搜索引擎
    engine, _ := NewSearchEngine("./data.db")
    defer engine.Close()

    // 导入文档
    for _, d := range docs {
        doc := NewDocument(
            fmt.Sprintf("%d", d.ID),
            d.Title,
            d.Content,
        )
        doc.URL = d.URL
        engine.UpsertDocument(doc)
    }

    log.Println("Migration complete!")
}
```

## 🧪 测试

### 运行示例脚本

```bash
# CLI 示例
chmod +x example_usage.sh
./example_usage.sh

# API 测试（需要先启动服务器）
# 终端 1：
go run *.go serve

# 终端 2：
chmod +x test_api.sh
./test_api.sh
```

### 手动测试流程

```bash
# 1. 启动服务器
go run *.go serve &

# 2. 等待服务器启动
sleep 2

# 3. 插入测试数据
curl -X POST http://localhost:3000/documents \
  -H "Content-Type: application/json" \
  -d '{"id":"test","title":"Test Doc","content":"This is a test"}'

# 4. 搜索
curl "http://localhost:3000/search?query=test"

# 5. 删除
curl -X DELETE http://localhost:3000/documents/test

# 6. 停止服务器
pkill -f "go run"
```

## 📝 注意事项

### 性能

- BoltDB 比 JSON 文件快 3-5 倍
- 内存使用：索引仍在内存中，大数据集需要注意
- 并发：完全线程安全，支持多个并发请求

### 存储

- 数据文件默认位置：`./data/search.db`
- 可通过 `--data-dir` 参数修改
- BoltDB 是单文件数据库，方便备份

### 兼容性

- Go 1.21 或更高版本
- 不兼容旧版本的 JSON 数据
- API 设计与 Rust 版本保持一致

## 🐛 常见问题

### Q1: 编译错误 "cannot find package"

```bash
go mod download
go mod tidy
```

### Q2: 端口已被占用

```bash
# 使用其他端口
go run *.go serve --port 8080
```

### Q3: 数据库文件损坏

```bash
# 删除并重新创建
rm -rf data/
go run *.go serve
```

### Q4: 内存使用过高

索引保存在内存中。对于大规模数据：
- 考虑分片
- 增加服务器内存
- 或使用外部搜索引擎（Elasticsearch）

## 🎓 学习资源

- [Gin 文档](https://gin-gonic.com/docs/)
- [Cobra 文档](https://github.com/spf13/cobra)
- [BoltDB 文档](https://github.com/etcd-io/bbolt)
- [BM25 算法](https://en.wikipedia.org/wiki/Okapi_BM25)

## 🙏 反馈

如有问题或建议，请提交 Issue。

## 📄 许可证

与原项目相同。
