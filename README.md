# Code Switch

集中管理 Claude Code & Codex 供应商

- 无需重启 cc & codex, 平滑切换不同供应商
- 支持多供应商自动降级, 保证使用体验
- 支持请求级别的用量统计, 花费多少清晰可见
- 支持 cc & codex Mcp Server 双平台管理
- 支持 Claude Skill 自动下载与安装, 内置 2 个流行的 skill 仓库
- 支持添加自定义 Skill 仓库

基于 [Wails 3](https://v3.wails.io)

## 实现原理

应用启动时会初始化 在本地 18100 端口创建一个 HTTP 代理服务器, 默认绑定 :18100

并自动更新 Claude Code、Codex 配置, 指向 http://127.0.0.1:18100 服务

代理内部只暴露兼容的关键端点：

- /v1/messages 转发到配置的 Claude 供应商
- /responses 转发到 Codex 供应商；

请求由 proxyHandler 动态挑选符合当前优先级与启用状态的 provider，并在失败时自动回退。

以上流程让 cli 看到的是一个固定的本地地址，而真实请求会被 Code Switch 透明地路由到你在应用里维护的供应商列表

## 下载

[macOS](https://github.com/daodao97/code-swtich/releases) | [windows](https://github.com/daodao97/code-swtich/releases) 


## 预览
![亮色主界面](resources/images/code-switch.png)
![暗色主界面](resources/images/code-swtich-dark.png)
![日志亮色](resources/images/code-switch-logs.png)
![日志暗色](resources/images/code-switch-logs-dark.png)

## 开发准备
- Go 1.24+
- Node.js 18+
- npm / pnpm / yarn
- Wails 3 CLI：`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`

## 开发运行
```bash
wails3 task dev
```

## 构建流程
1. 同步 build metadata：
   ```bash
   wails3 task common:update:build-assets
   ```
2. 打包 macOS `.app`：
   ```bash
   wails3 task package
   ```

### 交叉编译 Windows (macOS 环境)
1. 安装 `mingw-w64`：
   ```bash
   brew install mingw-w64
   ```
2. 运行 Windows 任务：
   ```bash
   env ARCH=amd64 wails3 task windows:build
   # 生成安装器
   env ARCH=amd64 wails3 task windows:package
   ```

## 发布

### 1. 准备发布说明
在项目根目录创建或更新 `RELEASE_NOTES.md` 文件，写入本次更新的内容（支持 Markdown）。

### 2. 运行发布脚本
脚本 `scripts/publish_r2.sh <version>` 将自动完成以下步骤：
1. 更新代码中的版本号
2. 构建前端资源
3. 打包 macOS (arm64 & amd64) 和 Windows (amd64) 应用
4. 上传所有安装包到 Cloudflare R2
5. 生成并上传 `latest.json`（包含版本信息和发布说明）

```bash
# 示例：发布 v1.0.0
./scripts/publish_r2.sh v1.0.0
```

发布完成后，用户打开应用时会自动检测到新版本并显示更新弹窗。

## 常见问题
- 若 `.app` 无法打开，先执行 `wails3 task common:update:build-assets` 后再构建。
- macOS 交叉编译需要终端拥有完全磁盘访问权限，否则 `~/Library/Caches/go-build` 会报 *operation not permitted*。
