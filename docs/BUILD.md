# Code Switch 构建指南

## 📦 环境构建

### QA 测试环境

```bash
# 方式 1: 使用构建脚本（推荐）
./build/qa.sh

# 方式 2: 手动构建
ENVIRONMENT=qa wails build -clean

# 方式 3: 开发模式
ENVIRONMENT=qa wails dev
```

**QA 环境配置**:
- 登录页面: `https://0011ai-qa.bihu.it/desktop-login`
- API 地址: `https://wildai-qa.bihu.it/`

### 生产环境

```bash
# 方式 1: 使用构建脚本（推荐）
./build/prod.sh

# 方式 2: 手动构建
ENVIRONMENT=production wails build -clean

# 方式 3: 使用默认值（不设置环境变量，默认为 production）
wails build -clean
```

**生产环境配置**:
- 登录页面: `https://0011.ai/desktop-login`
- API 地址: `https://api.wildai.net/`

---

## 🛠️ 开发模式

### 启动开发服务器

```bash
# QA 环境
ENVIRONMENT=qa wails dev

# 生产环境
ENVIRONMENT=production wails dev

# 默认（生产环境）
wails dev
```

### 热重载

开发模式下，前端代码修改会自动热重载，Go 代码修改需要重启。

---

## 📋 构建产物

### macOS

```
build/bin/
└── Code Switch.app/
```

**运行**:
```bash
open "build/bin/Code Switch.app"
```

### Windows

```
build/bin/
└── code-switch.exe
```

**运行**:
```bash
./build/bin/code-switch.exe
```

---

## ✅ 验证构建

### 检查环境配置

启动应用后，检查日志输出：

```
[INFO] Environment: qa
[INFO] Login URL: https://0011ai-qa.bihu.it/desktop-login
[INFO] API URL: https://wildai-qa.bihu.it/
```

或

```
[INFO] Environment: production
[INFO] Login URL: https://0011.ai/desktop-login
[INFO] API URL: https://api.wildai.net/
```

### 测试登录流程

1. 点击"登录"按钮
2. 应该打开正确环境的登录页面
3. 完成登录后返回应用

---

## 🔧 常见问题

### Q: 如何切换环境？

**A**: 重新构建应用，设置不同的 `ENVIRONMENT` 变量。

### Q: 构建时如何指定输出目录？

**A**: 修改 `wails.json`:
```json
{
  "outputDir": "build/bin"
}
```

### Q: 如何构建不同架构？

**A**:
```bash
# macOS Intel
ENVIRONMENT=production wails build -platform darwin/amd64

# macOS Apple Silicon
ENVIRONMENT=production wails build -platform darwin/arm64

# Windows 64位
ENVIRONMENT=production wails build -platform windows/amd64
```

### Q: 开发时想测试生产环境的 API？

**A**: 修改 `config/config.go`，添加开发环境配置，或直接在代码中临时修改 URL。

---

## 📦 发布检查清单

### QA 版本

- [ ] 使用 `./build/qa.sh` 构建
- [ ] 验证登录页面指向 QA 环境
- [ ] 测试完整登录流程
- [ ] 检查 API 调用是否指向 QA 环境
- [ ] 版本号标记为 `x.x.x-qa`

### 生产版本

- [ ] 使用 `./build/prod.sh` 构建
- [ ] 验证登录页面指向生产环境
- [ ] 测试完整登录流程
- [ ] 检查 API 调用是否指向生产环境
- [ ] 版本号使用正式版本号 `x.x.x`
- [ ] 代码签名（macOS/Windows）
- [ ] 创建安装包
- [ ] 上传到 GitHub Releases

---

## 🚀 CI/CD 集成

### GitHub Actions 示例

```yaml
name: Build

on:
  push:
    branches: [ main, develop ]

jobs:
  build-qa:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build QA
        run: |
          ENVIRONMENT=qa wails build -clean
      - name: Upload Artifact
        uses: actions/upload-artifact@v3
        with:
          name: code-switch-qa
          path: build/bin/

  build-prod:
    runs-on: macos-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build Production
        run: |
          ENVIRONMENT=production wails build -clean
      - name: Upload Artifact
        uses: actions/upload-artifact@v3
        with:
          name: code-switch-prod
          path: build/bin/
```

---

## 📚 相关文档

- [Wails 官方文档](https://wails.io/docs/)
- [构建选项](https://wails.io/docs/reference/cli#build)
- [Deep Link 登录功能](./DEEP_LINK_LOGIN.md)
