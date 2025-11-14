# OAuth 登录功能

## 📋 概述

基于 OAuth 2.0 + PKCE 标准实现的安全登录流程，参考 VSCode 的登录方案。用户在桌面应用点击"登录"按钮后，通过本地 HTTP 服务器 + Deep Link 的方式实现安全的认证流程。

## 🔄 新的登录流程

```
用户点击"登录"
    ↓
启动本地 HTTP 服务器 (127.0.0.1:52779-52879)
    ↓
生成安全参数 (nonce, PKCE, windowId)
    ↓
打开浏览器访问 Web 登录页 (带安全参数)
    ↓
用户在浏览器完成登录
    ↓
Web 端重定向到本地服务器 (带授权码)
    ↓
本地服务器验证 state 参数
    ↓
重定向到 Deep Link: codeswitch://auth?nonce=xxx&windowId=xxx&code=xxx
    ↓
桌面应用验证 nonce 并使用授权码交换 token
    ↓
保存 Token 到 ~/.code-switch/auth.json
    ↓
关闭本地服务器
    ↓
登录完成 ✅
```

## 🔒 安全机制

### 1. 本地 HTTP 服务器
- **临时端口**：随机选择 52779-52879 范围内的端口
- **生命周期**：登录完成后立即关闭
- **隔离性**：只监听 127.0.0.1，确保本地访问

### 2. PKCE (Proof Key for Code Exchange)
- **Code Challenge**：防止授权码拦截攻击
- **Code Verifier**：客户端证明身份的密钥
- **S256 方法**：使用 SHA256 哈希算法

### 3. Nonce 防护
- **防 CSRF**：每次登录使用唯一的随机数
- **时效性**：10 分钟有效期
- **一次性**：使用后立即标记为无效

### 4. 状态验证
- **WindowId**：标识发起登录的具体窗口
- **State 参数**：验证请求和响应的一致性
- **错误处理**：完善的错误检测和用户提示

## 📂 核心文件

### 后端 (Go)

1. **services/oauth_server.go** - 本地 OAuth 服务器
   - `Start()` - 启动本地 HTTP 服务器
   - `WaitForResult()` - 等待登录结果
   - `handleCallback()` - 处理 Web 端回调

2. **services/authservice.go** - 认证服务（暴露给前端）
   - `OpenLoginPage()` - 启动 OAuth 登录流程
   - `exchangeCodeForToken()` - 使用授权码交换 token
   - `GetAuthInfo()` - 获取当前认证信息
   - `Logout()` - 退出登录

3. **services/security.go** - 安全管理器
   - `GenerateNonce()` - 生成安全随机数
   - `GeneratePKCE()` - 生成 PKCE 参数
   - `ValidateNonce()` - 验证 nonce 参数

4. **services/deeplink.go** - Deep Link 处理
   - `HandleDeepLink()` - 处理 OAuth 和 Legacy 两种格式
   - `SaveAuthInfo()` - 保存认证信息到本地文件
   - `LoadAuthInfo()` - 加载认证信息

5. **main.go** - Deep Link 事件处理
   - 监听 `events.Common.ApplicationLaunchedWithUrl` 事件
   - 登录成功后发送 `auth:login-success` 事件到前端

6. **config/config.go** - 环境配置
   - Development: `http://localhost:3000/desktop-login`
   - QA: `https://0011ai-qa.bihu.it/desktop-login`
   - Production: `https://0011.ai/desktop-login`

7. **build/darwin/Info.plist** - macOS URL Scheme 注册
   - 注册 `codeswitch://` 协议

### 前端 (Vue 3)

1. **frontend/src/components/Main/Index.vue**
   - 主界面顶部的"登录"按钮
   - `handleLogin()` 函数调用 `AuthService.OpenLoginPage()`

2. **frontend/src/style.css**
   - `.login-button` 样式 - 蓝色按钮带文字和图标

3. **frontend/src/locales/**
   - 多语言支持（中文/英文）

## 🔧 登录 URL 参数

Web 登录页面会接收以下参数：

```http
GET /desktop-login?callback=http://127.0.0.1:52780/callback&state=xxx&window_id=xxx&code_challenge=xxx&code_challenge_method=S256
```

- **callback**: 本地服务器回调地址
- **state**: 防 CSRF 随机数
- **window_id**: 窗口标识符
- **code_challenge**: PKCE 码挑战
- **code_challenge_method**: PKCE 方法 (固定为 S256)

## 🎯 Web 端实现要求

Web 登录页面需要：

1. **接收参数**：解析上述 URL 参数
2. **完成登录**：正常的用户登录流程
3. **重定向回调**：登录成功后重定向到 callback 地址

```http
GET http://127.0.0.1:52780/callback?code=auth_code_here&state=xxx
```

4. **错误处理**：登录失败时包含错误参数

```http
GET http://127.0.0.1:52780/callback?error=access_denied&state=xxx
```

## 🛠️ 开发和构建

### 本地开发

```bash
# 默认使用 development 环境
wails3 dev

# 指定环境
ENVIRONMENT=qa wails3 dev
ENVIRONMENT=production wails3 dev
```

### 构建应用

```bash
# QA 环境
./build/qa.sh

# 生产环境
./build/prod.sh
```

## 🧪 测试

### 完整登录流程测试

1. **启动 Wails 应用**:
   ```bash
   wails3 dev
   ```

2. **点击登录按钮**:
   - 应用会在 52779-52879 范围内启动本地 HTTP 服务器
   - 打开浏览器访问登录页面（带安全参数）

3. **在浏览器完成登录**:
   - 输入用户名密码
   - 点击登录

4. **验证流程**:
   - 浏览器应该重定向到本地服务器
   - 本地服务器应该重定向到 Deep Link
   - 应用应该接收到 Deep Link 并完成登录

5. **检查结果**:
   - Token 应保存到 `~/.code-switch/auth.json`
   - 前端应该显示登录状态

### 手动测试 OAuth 回调

```bash
# macOS
open "codeswitch://auth?nonce=test123&windowId=456&code=sample_auth_code"

# Windows
Start-Process "codeswitch://auth?nonce=test123&windowId=456&code=sample_auth_code"
```

### 查看认证信息

```bash
# 查看保存的认证信息
cat ~/.code-switch/auth.json

# 手动清除
rm ~/.code-switch/auth.json
```

## 🔍 调试

### 查看本地服务器状态

```bash
# 检查端口是否被占用
lsof -i :52779-52879

# 测试本地服务器健康检查
curl http://127.0.0.1:52779/health
```

### 调试日志

应用会输出详细的调试日志：
- OAuth 服务器启动/停止
- 登录 URL 参数
- 回调处理结果
- Token 交换过程

## 🆚 vs 原始方案对比

| 特性 | 原始 Deep Link | 新 OAuth 方案 |
|------|---------------|---------------|
| **Token 传递** | URL 明文参数 | 授权码 + 安全交换 |
| **安全性** | ⚠️ 低 | ✅ 高 (OAuth 2.0) |
| **防 CSRF** | ❌ 无 | ✅ Nonce 验证 |
| **防拦截** | ❌ 无 | ✅ PKCE |
| **用户体验** | 直接跳转 | 自动重定向 |
| **标准化** | 自定义 | ✅ 行业标准 |
| **兼容性** | ✅ 向后兼容 | ✅ 支持新旧格式 |

## 📝 实现特点

✅ **安全第一** - 符合 OAuth 2.0 + PKCE 安全标准
✅ **用户体验** - 无感知的自动重定向
✅ **向后兼容** - 支持旧的 Direct Token 格式
✅ **多环境支持** - Development/QA/Production
✅ **错误处理** - 完善的错误检测和用户提示
✅ **资源管理** - 本地服务器用完即关闭
✅ **跨平台** - 支持 macOS/Windows/Linux

---

**创建时间**: 2025-01-14
**版本**: v2.0 (OAuth 方案)
**状态**: ✅ 开发完成