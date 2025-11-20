#!/bin/bash

# 修复数据库问题的脚本

echo "🔧 修复 Code Switch 数据库"
echo "================================"
echo ""

DB_PATH=~/.code-switch/app.db

# 检查数据库文件
if [ ! -f "$DB_PATH" ]; then
    echo "❌ 数据库文件不存在，将创建新数据库"
    mkdir -p ~/.code-switch
    touch "$DB_PATH"
fi

# 检查数据库大小
DB_SIZE=$(stat -f%z "$DB_PATH" 2>/dev/null || stat -c%s "$DB_PATH" 2>/dev/null)

if [ "$DB_SIZE" -eq 0 ]; then
    echo "⚠️  数据库文件为空 (0 字节)"
    echo ""
fi

# 检查表是否存在
echo "检查数据库表..."
TABLES=$(sqlite3 "$DB_PATH" ".tables" 2>&1)

if [[ "$TABLES" == *"request_log"* ]]; then
    echo "✅ request_log 表已存在"
    
    # 检查记录数
    COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM request_log;" 2>/dev/null)
    echo "   记录数: $COUNT"
else
    echo "❌ request_log 表不存在，正在创建..."
    
    # 创建表
    sqlite3 "$DB_PATH" "CREATE TABLE IF NOT EXISTS request_log (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        platform TEXT,
        model TEXT,
        provider TEXT,
        http_code INTEGER,
        input_tokens INTEGER,
        output_tokens INTEGER,
        cache_create_tokens INTEGER,
        cache_read_tokens INTEGER,
        reasoning_tokens INTEGER,
        is_stream INTEGER DEFAULT 0,
        duration_sec REAL DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );" 2>&1
    
    if [ $? -eq 0 ]; then
        echo "✅ 表创建成功"
    else
        echo "❌ 表创建失败"
        exit 1
    fi
fi

echo ""
echo "================================"
echo "✅ 数据库修复完成"
echo "================================"
echo ""
echo "数据库路径: $DB_PATH"
echo "数据库大小: $(du -h "$DB_PATH" | cut -f1)"
echo ""
echo "现在可以："
echo "1. 刷新应用界面 (Cmd+R)"
echo "2. 或重启应用"
echo ""
