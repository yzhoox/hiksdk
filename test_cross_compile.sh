#!/bin/bash

# 跨平台编译测试脚�?echo "==================================="
echo "跨平台编译测�?
echo "==================================="

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # 无颜�?
# 确保 CGO 已启�?export CGO_ENABLED=1

echo ""
echo "1. 测试当前平台编译..."
echo "-----------------------------------"
if go build -v ./core; then
    echo -e "${GREEN}�?当前平台编译成功${NC}"
else
    echo -e "${RED}�?当前平台编译失败${NC}"
    exit 1
fi

echo ""
echo "2. 测试 Linux 平台编译（仅语法检查）..."
echo "-----------------------------------"
if GOOS=linux GOARCH=amd64 go build -o /dev/null ./core 2>&1 | grep -q "error"; then
    echo -e "${RED}�?Linux 平台存在编译错误${NC}"
    GOOS=linux GOARCH=amd64 go build -o /dev/null ./core
else
    echo -e "${GREEN}�?Linux 平台语法检查通过${NC}"
fi

echo ""
echo "3. 测试 Windows 平台编译（仅语法检查）..."
echo "-----------------------------------"
if GOOS=windows GOARCH=amd64 go build -o /dev/null ./core 2>&1 | grep -q "error"; then
    echo -e "${RED}�?Windows 平台存在编译错误${NC}"
    GOOS=windows GOARCH=amd64 go build -o /dev/null ./core
else
    echo -e "${GREEN}�?Windows 平台语法检查通过${NC}"
fi

echo ""
echo "==================================="
echo "测试完成"
echo "==================================="
