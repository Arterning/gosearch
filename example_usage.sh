#!/bin/bash

# SimpleFTS V2 使用示例脚本

echo "=== SimpleFTS V2 使用示例 ==="
echo ""

# 设置颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}1. 插入一些测试文档${NC}"
echo ""

echo "插入文档 1..."
go run *.go insert \
  --id "1" \
  --title "Go Programming Language" \
  --content "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software. Go is statically typed and compiled." \
  --url "https://golang.org"

echo ""
echo "插入文档 2..."
go run *.go insert \
  --id "2" \
  --title "Rust Programming Language" \
  --content "Rust is a systems programming language that runs blazingly fast, prevents segfaults, and guarantees thread safety." \
  --url "https://www.rust-lang.org"

echo ""
echo "插入文档 3..."
go run *.go insert \
  --id "3" \
  --title "Python Programming" \
  --content "Python is a programming language that lets you work quickly and integrate systems more effectively. Python is interpreted and dynamically typed." \
  --url "https://www.python.org"

echo ""
echo "插入文档 4..."
go run *.go insert \
  --id "4" \
  --title "JavaScript Introduction" \
  --content "JavaScript is the programming language of the Web. JavaScript is easy to learn and can run on both client and server." \
  --url "https://www.javascript.com"

echo ""
echo -e "${GREEN}✓ 已插入 4 个文档${NC}"
echo ""

sleep 1

echo -e "${BLUE}2. 搜索 'programming language'${NC}"
go run *.go search --query "programming language" --limit 10

echo ""
sleep 1

echo -e "${BLUE}3. 搜索 'fast efficient' (AND 模式)${NC}"
go run *.go search --query "fast efficient" --mode and

echo ""
sleep 1

echo -e "${BLUE}4. 搜索 'rust python' (OR 模式)${NC}"
go run *.go search --query "rust python" --mode or

echo ""
sleep 1

echo -e "${BLUE}5. 获取文档 1${NC}"
go run *.go get --id "1"

echo ""
sleep 1

echo -e "${BLUE}6. 查看索引统计${NC}"
go run *.go stats

echo ""
echo -e "${BLUE}7. 删除文档 4${NC}"
go run *.go delete --id "4"

echo ""
echo -e "${GREEN}✓ 示例完成！${NC}"
echo ""
echo "现在可以尝试启动 HTTP 服务器："
echo "  go run *.go serve"
echo ""
echo "然后使用 curl 测试 API："
echo "  curl http://localhost:3000/search?query=programming"
