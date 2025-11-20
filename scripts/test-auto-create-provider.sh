#!/bin/bash

# 测试自动创建 0011 供应商功能的脚本

set -e

BACKUP_DIR="$HOME/.code-switch/backup-$(date +%Y%m%d-%H%M%S)"
CONFIG_DIR="$HOME/.code-switch"

echo "🧪 开始测试自动创建 0011 供应商功能"
echo ""

# 检查配置目录是否存在
if [ ! -d "$CONFIG_DIR" ]; then
    echo "❌ 配置目录不存在: $CONFIG_DIR"
    exit 1
fi

# 创建备份目录
echo "📦 创建备份目录: $BACKUP_DIR"
mkdir -p "$BACKUP_DIR"

# 备份现有配置
if [ -f "$CONFIG_DIR/claude-code.json" ]; then
    echo "📋 备份 claude-code.json"
    cp "$CONFIG_DIR/claude-code.json" "$BACKUP_DIR/"
fi

if [ -f "$CONFIG_DIR/codex.json" ]; then
    echo "📋 备份 codex.json"
    cp "$CONFIG_DIR/codex.json" "$BACKUP_DIR/"
fi

echo ""
echo "请选择测试场景："
echo "1) 测试从零创建（删除所有配置）"
echo "2) 测试更新 API Key（保留配置但清空 0011 的 API Key）"
echo "3) 恢复备份"
read -p "请输入选项 (1/2/3): " choice

case $choice in
    1)
        echo ""
        echo "🗑️  删除现有配置文件..."
        rm -f "$CONFIG_DIR/claude-code.json"
        rm -f "$CONFIG_DIR/codex.json"
        echo "✅ 配置已清空"
        echo ""
        echo "📝 现在请重启应用并登录，观察是否自动创建 0011 卡片"
        ;;
    2)
        echo ""
        echo "🔧 清空 0011 供应商的 API Key..."
        
        # 使用 jq 修改 JSON（如果安装了 jq）
        if command -v jq &> /dev/null; then
            if [ -f "$CONFIG_DIR/claude-code.json" ]; then
                jq '(.providers[] | select(.id == 100) | .apiKey) = ""' "$CONFIG_DIR/claude-code.json" > "$CONFIG_DIR/claude-code.json.tmp"
                mv "$CONFIG_DIR/claude-code.json.tmp" "$CONFIG_DIR/claude-code.json"
                echo "✅ 已清空 Claude Code 0011 的 API Key"
            fi
            
            if [ -f "$CONFIG_DIR/codex.json" ]; then
                jq '(.providers[] | select(.id == 200) | .apiKey) = ""' "$CONFIG_DIR/codex.json" > "$CONFIG_DIR/codex.json.tmp"
                mv "$CONFIG_DIR/codex.json.tmp" "$CONFIG_DIR/codex.json"
                echo "✅ 已清空 Codex 0011 的 API Key"
            fi
        else
            echo "⚠️  未安装 jq，请手动编辑配置文件："
            echo "   $CONFIG_DIR/claude-code.json"
            echo "   $CONFIG_DIR/codex.json"
            echo ""
            echo "将 ID 为 100 和 200 的卡片的 apiKey 字段改为空字符串"
        fi
        
        echo ""
        echo "📝 现在请重启应用并登录，观察是否自动更新 API Key"
        ;;
    3)
        echo ""
        echo "🔄 恢复备份..."
        if [ -f "$BACKUP_DIR/claude-code.json" ]; then
            cp "$BACKUP_DIR/claude-code.json" "$CONFIG_DIR/"
            echo "✅ 已恢复 claude-code.json"
        fi
        if [ -f "$BACKUP_DIR/codex.json" ]; then
            cp "$BACKUP_DIR/codex.json" "$CONFIG_DIR/"
            echo "✅ 已恢复 codex.json"
        fi
        ;;
    *)
        echo "❌ 无效选项"
        exit 1
        ;;
esac

echo ""
echo "📂 备份位置: $BACKUP_DIR"
echo ""
