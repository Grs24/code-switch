#!/bin/bash

# 测试代理的流式响应

echo "🧪 测试代理流式响应"
echo "================================"
echo ""

# 获取 API Key
API_KEY=$(cat ~/.code-switch/claude-code.json | jq -r '.providers[] | select(.id == 100) | .apiKey')

if [ -z "$API_KEY" ] || [ "$API_KEY" = "null" ]; then
    echo "❌ 无法获取 API Key"
    exit 1
fi

echo "✅ API Key: ${API_KEY:0:10}..."
echo ""

# 测试流式请求
echo "发送测试请求到代理..."
echo ""

curl -N -X POST http://127.0.0.1:18100/v1/messages \
  -H "Content-Type: application/json" \
  -H "anthropic-version: 2023-06-01" \
  -H "x-api-key: $API_KEY" \
  -d '{
    "model": "claude-sonnet-4-20250514",
    "max_tokens": 100,
    "stream": true,
    "messages": [
      {
        "role": "user",
        "content": "Say hello in one word"
      }
    ]
  }' 2>&1 | tee /tmp/proxy-response.txt

echo ""
echo ""
echo "================================"
echo "响应已保存到: /tmp/proxy-response.txt"
echo ""
echo "检查响应："
echo "- 是否包含 'data:' 行？"
echo "- 是否以 'data: [DONE]' 结束？"
echo "- 连接是否正确关闭？"
echo ""
