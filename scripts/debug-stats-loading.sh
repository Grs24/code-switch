#!/bin/bash

# 排查统计数据加载问题的脚本

echo "🔍 排查统计数据加载问题"
echo "================================"
echo ""

# 检查数据库文件
echo "1️⃣ 检查数据库文件"
echo "--------------------------------"
if [ -f ~/.code-switch/app.db ]; then
    DB_SIZE=$(du -h ~/.code-switch/app.db | cut -f1)
    echo "✅ 数据库文件存在"
    echo "   大小: $DB_SIZE"
    echo ""
    
    # 检查数据库中的记录数
    echo "   查询记录数..."
    sqlite3 ~/.code-switch/app.db "SELECT COUNT(*) FROM request_log;" 2>/dev/null && echo "" || echo "   ⚠️  无法查询数据库"
else
    echo "❌ 数据库文件不存在"
    echo ""
fi

# 检查代理服务
echo "2️⃣ 检查代理服务"
echo "--------------------------------"
if lsof -i:18100 > /dev/null 2>&1; then
    echo "✅ 代理服务正在运行 (端口 18100)"
    lsof -i:18100 | grep LISTEN
else
    echo "❌ 代理服务未运行"
fi
echo ""

# 检查最近的日志
echo "3️⃣ 检查最近的请求日志"
echo "--------------------------------"
if [ -f ~/.code-switch/app.db ]; then
    echo "最近 5 条请求记录："
    sqlite3 ~/.code-switch/app.db "SELECT datetime(created_at, 'unixepoch', 'localtime') as time, platform, provider, model, http_code FROM request_log ORDER BY created_at DESC LIMIT 5;" 2>/dev/null || echo "   ⚠️  无法查询"
else
    echo "   数据库文件不存在"
fi
echo ""

# 测试统计查询性能
echo "4️⃣ 测试统计查询性能"
echo "--------------------------------"
if [ -f ~/.code-switch/app.db ]; then
    echo "执行统计查询..."
    time sqlite3 ~/.code-switch/app.db "
    SELECT 
        provider,
        COUNT(*) as total_requests,
        SUM(input_tokens) as input_tokens,
        SUM(output_tokens) as output_tokens
    FROM request_log 
    WHERE created_at >= strftime('%s', 'now', 'start of day', '-1 day')
    GROUP BY provider;
    " 2>/dev/null || echo "   ⚠️  查询失败"
else
    echo "   数据库文件不存在"
fi
echo ""

# 检查数据库锁
echo "5️⃣ 检查数据库锁"
echo "--------------------------------"
if lsof ~/.code-switch/app.db 2>/dev/null; then
    echo "✅ 数据库被以下进程使用："
    lsof ~/.code-switch/app.db
else
    echo "   没有进程锁定数据库"
fi
echo ""

# 检查网络连接
echo "6️⃣ 检查网络连接"
echo "--------------------------------"
echo "测试本地连接..."
if curl -s http://127.0.0.1:18100 > /dev/null 2>&1; then
    echo "✅ 本地代理可访问"
else
    echo "⚠️  本地代理无响应"
fi
echo ""

# 建议
echo "================================"
echo "💡 排查建议"
echo "================================"
echo ""
echo "如果看到 '刷新中' 状态一直不消失："
echo ""
echo "1. 检查浏览器控制台错误"
echo "   - 按 Cmd+Option+I 打开开发者工具"
echo "   - 查看 Console 标签是否有错误"
echo ""
echo "2. 检查网络请求"
echo "   - 在开发者工具的 Network 标签"
echo "   - 查找 'ProviderDailyStats' 请求"
echo "   - 查看请求状态和响应时间"
echo ""
echo "3. 如果数据库很大，考虑清理旧数据"
echo "   sqlite3 ~/.code-switch/app.db \"DELETE FROM request_log WHERE created_at < strftime('%s', 'now', '-30 day');\""
echo ""
echo "4. 如果数据库被锁定，重启应用"
echo "   pkill -f CodeSwitch"
echo "   ./start-desktop.sh"
echo ""
echo "5. 查看后端日志"
echo "   - 在运行 start-desktop.sh 的终端查看输出"
echo "   - 查找是否有数据库错误"
echo ""
