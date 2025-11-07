#!/bin/bash

# SimpleFTS V2 API 测试脚本
# 使用前请先启动服务器：go run *.go serve

API_URL="http://localhost:3000"

echo "=== SimpleFTS V2 API 测试 ==="
echo ""
echo "确保服务器已启动：go run *.go serve"
echo ""

sleep 2

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}1. 健康检查${NC}"
curl -s ${API_URL}/health | jq
echo ""

sleep 1

echo -e "${BLUE}2. 插入文档 1${NC}"
curl -s -X POST ${API_URL}/documents \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "title": "Go Programming Language",
    "content": "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
    "url": "https://golang.org"
  }' | jq
echo ""

sleep 1

echo -e "${BLUE}3. 批量插入文档${NC}"
curl -s -X POST ${API_URL}/documents/batch \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "id": "2",
        "title": "Rust Programming",
        "content": "Rust is a systems programming language that runs blazingly fast and prevents segfaults."
      },
      {
        "id": "3",
        "title": "Python Programming",
        "content": "Python is a high-level programming language that lets you work quickly."
      }
    ]
  }' | jq
echo ""

sleep 1

echo -e "${BLUE}4. 搜索 'programming language'${NC}"
curl -s "${API_URL}/search?query=programming+language&limit=5" | jq
echo ""

sleep 1

echo -e "${BLUE}5. 搜索 'fast' (带 BM25 排序)${NC}"
curl -s "${API_URL}/search?query=fast&ranked=true" | jq
echo ""

sleep 1

echo -e "${BLUE}6. OR 搜索 'rust python'${NC}"
curl -s "${API_URL}/search?query=rust+python&mode=or" | jq
echo ""

sleep 1

echo -e "${BLUE}7. 获取文档 1${NC}"
curl -s ${API_URL}/documents/1 | jq
echo ""

sleep 1

echo -e "${BLUE}8. 更新文档 1${NC}"
curl -s -X PUT ${API_URL}/documents/1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "title": "Go: A Modern Programming Language",
    "content": "Go is a statically typed, compiled programming language designed at Google."
  }' | jq
echo ""

sleep 1

echo -e "${BLUE}9. 再次获取文档 1（验证更新）${NC}"
curl -s ${API_URL}/documents/1 | jq
echo ""

sleep 1

echo -e "${BLUE}10. 获取统计信息${NC}"
curl -s ${API_URL}/stats | jq
echo ""

sleep 1

echo -e "${BLUE}11. 删除文档 3${NC}"
curl -s -X DELETE ${API_URL}/documents/3 | jq
echo ""

sleep 1

echo -e "${BLUE}12. 再次搜索验证删除${NC}"
curl -s "${API_URL}/search?query=python" | jq
echo ""

echo -e "${GREEN}=== API 测试完成 ===${NC}"
echo ""
echo -e "${YELLOW}提示：如果看到 'command not found: jq'，请安装 jq 工具${NC}"
echo "或者移除脚本中的 '| jq' 查看原始 JSON"
