# SimpleFTS V2 - Enhanced Full-Text Search Engine

原始 Go 版本的增强版本，新增以下功能：

## 🆕 新功能

- ✅ **完整的 CRUD API** - 插入、更新、删除、查询文档
- ✅ **BM25 相关性排序** - 智能的文档相关性评分
- ✅ **HTTP REST API** - 基于 Gin 框架的 Web API
- ✅ **BoltDB 持久化** - 替代 JSON 文件，更高效
- ✅ **增量更新** - 实时插入和删除文档
- ✅ **分页支持** - 灵活的结果分页
- ✅ **AND/OR 搜索模式** - 支持多种查询模式
- ✅ **CLI 子命令** - 使用 Cobra 实现完整的命令行工具

## 📦 安装依赖

```bash
go mod download
```

## 🚀 使用方法


### 直接运行

```bash
go run main.go document.go index.go storage.go ranking.go engine.go api.go tokenizer.go filter.go [command]
```

## 📚 命令行使用

### 启动 HTTP 服务器

```bash
# 默认启动（127.0.0.1:3000）
go run *.go serve

# 自定义主机和端口
go run *.go serve --host 0.0.0.0 --port 8080

# 指定数据目录
go run *.go serve --data-dir ./my_data.db
```

### 插入文档

```bash
go run *.go insert \
  --id "doc1" \
  --title "Go Programming" \
  --content "Go is a simple and efficient programming language" \
  --url "https://golang.org"
```

### 搜索文档

```bash
# 基本搜索
go run *.go search --query "programming language"

# 指定返回数量
go run *.go search --query "go" --limit 5

# OR 搜索
go run *.go search --query "go rust python" --mode or

# 不使用排序
go run *.go search --query "programming" --ranked=false
```

### 获取文档

```bash
go run *.go get --id "doc1"
```

### 删除文档

```bash
go run *.go delete --id "doc1"
```

### 查看统计

```bash
go run *.go stats
```

## 🌐 HTTP API 使用

启动服务器后：

### 1. 健康检查

```bash
curl http://localhost:3000/health
```

### 2. 插入单个文档

```bash
curl -X POST http://localhost:3000/documents \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "title": "Go Programming Language",
    "content": "Go is a statically typed, compiled programming language",
    "url": "https://golang.org"
  }'
```

### 3. 批量插入文档

```bash
curl -X POST http://localhost:3000/documents/batch \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "id": "2",
        "title": "Rust Programming",
        "content": "Rust is a systems programming language"
      },
      {
        "id": "3",
        "title": "Python Programming",
        "content": "Python is an interpreted high-level language"
      }
    ]
  }'
```

### 4. 搜索文档

```bash
# 基本搜索
curl "http://localhost:3000/search?query=programming+language"

# 带参数的搜索
curl "http://localhost:3000/search?query=rust&limit=5&offset=0&ranked=true&mode=and"
```

**查询参数：**
- `query` - 搜索查询（必需）
- `limit` - 返回结果数量（默认: 10）
- `offset` - 分页偏移量（默认: 0）
- `ranked` - 是否使用 BM25 排序（默认: true）
- `mode` - 搜索模式：`and`（全匹配）或 `or`（任意匹配，默认: and）

### 5. 获取文档

```bash
curl http://localhost:3000/documents/1
```

### 6. 更新文档

```bash
curl -X PUT http://localhost:3000/documents/1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "title": "Updated Title",
    "content": "Updated content"
  }'
```

### 7. 删除文档

```bash
curl -X DELETE http://localhost:3000/documents/1
```

### 8. 获取统计信息

```bash
curl http://localhost:3000/stats
```

## 🏗️ 架构说明

### 新增模块

- `document.go` - 文档结构定义
- `index.go` - 改进的倒排索引（支持 CRUD）
- `storage.go` - BoltDB 持久化层
- `ranking.go` - BM25 排序算法
- `engine.go` - 搜索引擎核心
- `api.go` - Gin HTTP API
- `main.go` - 主程序（支持 CLI 和 Server）

### 保留模块

- `tokenizer.go` - 分词器
- `filter.go` - 文本过滤器（lowercase, stopword, stemmer）

## 📊 性能对比

| 特性 | 原版本 | V2 版本 |
|------|--------|---------|
| 存储方式 | JSON 文件 | BoltDB |
| 更新方式 | 全量重写 | 增量更新 |
| API 接口 | 无 | HTTP REST API |
| 相关性排序 | 无 | BM25 算法 |
| 文档操作 | 只读 | 完整 CRUD |
| 分页支持 | 无 | 支持 |
| 并发安全 | 部分 | 完全支持 |

## 🔄 迁移指南

从原版本迁移：

1. 旧数据不兼容，需要重新导入
2. 原有的 `document.go`、`index.go`、`main.go` 已被替换
3. 如需保留旧功能，请备份原文件

## 🐛 测试

```bash
# 运行测试
go test ./...

# 快速测试流程
# 1. 启动服务器
go run *.go serve &

# 2. 插入测试文档
curl -X POST http://localhost:3000/documents \
  -H "Content-Type: application/json" \
  -d '{"id":"test1","title":"Test","content":"This is a test document"}'

# 3. 搜索
curl "http://localhost:3000/search?query=test"

# 4. 删除
curl -X DELETE http://localhost:3000/documents/test1
```

## 📝 注意事项

1. **数据存储位置**：默认为 `./data/search.db`，可通过 `--data-dir` 修改
2. **BoltDB 文件**：单个文件数据库，方便备份
3. **并发安全**：所有操作都是线程安全的
4. **内存使用**：索引保存在内存中，大规模数据需注意内存使用

## 🎯 下一步

- [ ] 添加中文分词支持
- [ ] 实现模糊搜索
- [ ] 添加搜索高亮
- [ ] 支持多字段搜索
- [ ] 添加搜索建议

## 📄 许可证

与原项目相同
