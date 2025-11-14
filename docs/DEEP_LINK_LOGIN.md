# Deep Link 登录功能 (已弃用)

## ⚠️ 重要说明

此文档描述的是**旧的 Direct Token 方案**，已被新的 **OAuth 2.0 方案**替代。

📖 **请使用新的 OAuth 方案**：见 [OAuth 登录功能](./OAUTH_LOGIN.md)

## 📋 历史概述

## 🔄 登录流程

```
用户点击"登录"
    ↓
打开浏览器访问 Web 登录页
    ↓
用户在浏览器完成登录
    ↓
浏览器重定向: codeswitch://auth?token=xxx&user_id=xxx&email=xxx
    ↓
桌面应用接收 Deep Link 回调
    ↓
保存 Token 到 ~/.code-switch/auth.json
    ↓
登录完成 ✅
```

## 📂 核心文件

### 后端 (Go)

1. **services/deeplink.go** - Deep Link URL 解析和 Token 存储
   - `HandleDeepLink()` - 处理 Deep Link 回调
   - `saveAuthInfo()` - 保存认证信息到本地文件
   - `LoadAuthInfo()` - 加载认证信息
   - `ClearAuthInfo()` - 清除认证信息

2. **services/authservice.go** - 认证服务（暴露给前端）
   - `OpenLoginPage()` - 在浏览器中打开登录页面
   - `GetAuthInfo()` - 获取当前认证信息
   - `Logout()` - 退出登录

3. **main.go** - Deep Link 事件处理
   - 监听 `events.Common.ApplicationLaunchedWithUrl` 事件
   - 登录成功后发送 `auth:login-success` 事件到前端

4. **config/config.go** - 环境配置
   - Development: `http://localhost:3000/desktop-login`
   - QA: `https://0011ai-qa.bihu.it/desktop-login`
   - Production: `https://0011.ai/desktop-login`

5. **build/darwin/Info.plist** - macOS URL Scheme 注册
   - 注册 `codeswitch://` 协议

### 前端 (Vue 3)

1. **frontend/src/components/Main/Index.vue**
   - 主界面顶部的"登录"按钮
   - `handleLogin()` 函数直接调用 `AuthService.OpenLoginPage()`

2. **frontend/src/style.css**
   - `.login-button` 样式 - 蓝色按钮带文字和图标

3. **frontend/src/locales/**
   - 多语言支持（中文/英文）

## 🛠️ 开发和构建

### 本地开发

```bash
# 默认使用 development 环境（http://localhost:3000）
wails3 dev

# 或指定环境
ENVIRONMENT=qa wails3 dev
ENVIRONMENT=production wails3 dev
```

### 构建应用

```bash
# QA 环境
./build/qa.sh
# 或
ENVIRONMENT=qa wails3 build

# 生产环境
./build/prod.sh
# 或
ENVIRONMENT=production wails3 build
```

## 🧪 测试

### 本地测试

1. **启动 Web 开发服务器**:
   ```bash
   cd /Users/rongs/Code/0011ai
   npm run start
   ```

2. **启动 Wails 应用**:
   ```bash
   cd /Users/rongs/Code/test/code-switch-wails
   wails3 dev
   ```

3. **测试流程**:
   - 点击应用右上角的"登录"按钮
   - 浏览器应该打开 `http://localhost:3000/desktop-login`
   - 完成登录后，应该自动返回应用

### 手动测试 Deep Link

```bash
# macOS
open "codeswitch://auth?token=test123&user_id=456&email=test@example.com"

# Windows
Start-Process "codeswitch://auth?token=test123&user_id=456&email=test@example.com"
```

**预期结果**:
- 如果应用正在运行，控制台会显示收到 Deep Link
- Token 保存到 `~/.code-switch/auth.json`

### 查看认证信息

```bash
# 查看保存的认证信息
cat ~/.code-switch/auth.json

# 手动清除
rm ~/.code-switch/auth.json
```

## 🔍 调试

### 查看 Deep Link 是否注册

```bash
# macOS
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep codeswitch

# Windows
reg query HKEY_CLASSES_ROOT\codeswitch
```

### 常见问题

**Q: 浏览器没有打开？**
- 检查控制台是否有错误
- 确认 AuthService 是否正确绑定

**Q: Deep Link 没有被接收？**
- 确认应用正在运行
- 检查 Info.plist 配置是否正确
- 重新构建应用

**Q: Token 没有保存？**
- 检查 `~/.code-switch/` 目录权限
- 查看控制台错误信息

## 🚀 部署清单

### QA 版本
- [ ] 使用 `./build/qa.sh` 构建
- [ ] 验证登录页面指向 `https://0011ai-qa.bihu.it/desktop-login`
- [ ] 测试完整登录流程
- [ ] 测试 Deep Link 回调

### 生产版本
- [ ] 使用 `./build/prod.sh` 构建
- [ ] 验证登录页面指向 `https://0011.ai/desktop-login`
- [ ] 测试完整登录流程
- [ ] 测试 Deep Link 回调
- [ ] 代码签名（macOS/Windows）
- [ ] 创建安装包
- [ ] 上传到 GitHub Releases

## 📝 实现特点

✅ **简洁高效** - 无额外页面和状态管理，直接调用服务
✅ **零后端改动** - 完全复用现有 Web 登录页面
✅ **多环境支持** - Development/QA/Production 自动切换
✅ **安全可靠** - Token 存储权限 0600，仅所有者可读写
✅ **跨平台** - 支持 macOS/Windows/Linux

---

**创建时间**: 2025-01-13
**版本**: v1.0
**状态**: ✅ 开发完成
