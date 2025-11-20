#!/bin/bash

# 测试 SetUserInfo 方法的脚本

echo "🧪 测试 SetUserInfo 调用"
echo ""
echo "请按照以下步骤操作："
echo ""
echo "1. 确保应用正在运行（wails3 task dev）"
echo "2. 打开浏览器开发者工具（在应用中按 Cmd+Option+I）"
echo "3. 在控制台中执行以下代码："
echo ""
echo "================================================"
cat << 'EOF'
// 测试调用 SetUserInfo
const testUserInfo = {
  token: "test-token-123456",
  userId: "test-user-id",
  email: "test@example.com",
  mobile: "13800138000",
  userInfo: {}
};

// 导入 Call 方法
import('@wailsio/runtime').then(({ Call }) => {
  Call.ByName('AuthService.SetUserInfo', JSON.stringify(testUserInfo))
    .then(() => {
      console.log('✅ SetUserInfo 调用成功');
    })
    .catch((error) => {
      console.error('❌ SetUserInfo 调用失败:', error);
    });
});
EOF
echo "================================================"
echo ""
echo "4. 查看终端输出，应该能看到 [AuthService] 的日志"
echo ""
