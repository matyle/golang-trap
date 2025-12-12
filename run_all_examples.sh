#!/bin/bash

# 运行所有 Go 陷阱示例

echo "=========================================="
echo "Go 语言陷阱示例 - 运行所有示例"
echo "=========================================="
echo ""

EXAMPLES_DIR="examples"

# 检查 examples 目录是否存在
if [ ! -d "$EXAMPLES_DIR" ]; then
    echo "错误: examples 目录不存在"
    exit 1
fi

# 获取所有 .go 文件
files=$(find "$EXAMPLES_DIR" -name "*.go" | sort)

if [ -z "$files" ]; then
    echo "错误: 没有找到示例文件"
    exit 1
fi

# 运行每个示例
for file in $files; do
    echo "----------------------------------------"
    echo "运行: $file"
    echo "----------------------------------------"
    go run "$file"
    echo ""
    sleep 0.5  # 短暂延迟，便于观察输出
done

echo "=========================================="
echo "所有示例运行完成"
echo "=========================================="

