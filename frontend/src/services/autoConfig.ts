import { Call } from '@wailsio/runtime'

const SERVICE_NAME = 'codeswitch/services.AutoConfigService'

export interface AutoConfigResult {
  success: boolean
  claudeStatus: string
  codexStatus: string
  errors: string[]
  warnings: string[]
  os: string
  configPaths: Record<string, string[]>
  verificationCommands: Record<string, string[]>
}

export interface SystemInfo {
  os: string
  arch: string
  shell: string
}

export interface DependencyCheckResult {
  passed: boolean
  nodejs: NodeJSInfo
  errors: string[]
  warnings: string[]
}

export interface NodeJSInfo {
  installed: boolean
  version: string
  meetsRequirement: boolean
  minVersion: string
}

/**
 * 自动配置 Claude Code 和 Codex
 */
export async function ConfigureAll(apiKey: string, baseURL: string = ''): Promise<AutoConfigResult> {
  return Call.ByName(`${SERVICE_NAME}.ConfigureAll`, apiKey, baseURL)
}

/**
 * 仅配置 Claude Code
 */
export async function ConfigureClaudeCodeOnly(apiKey: string, baseURL: string = ''): Promise<AutoConfigResult> {
  return Call.ByName(`${SERVICE_NAME}.ConfigureClaudeCodeOnly`, apiKey, baseURL)
}

/**
 * 仅配置 Codex
 */
export async function ConfigureCodexOnly(apiKey: string, baseURL: string = ''): Promise<AutoConfigResult> {
  return Call.ByName(`${SERVICE_NAME}.ConfigureCodexOnly`, apiKey, baseURL)
}

/**
 * 验证配置是否正确
 */
export async function VerifyConfiguration(): Promise<AutoConfigResult> {
  return Call.ByName(`${SERVICE_NAME}.VerifyConfiguration`)
}

/**
 * 获取系统信息
 */
export async function GetSystemInfo(): Promise<SystemInfo> {
  return Call.ByName(`${SERVICE_NAME}.GetSystemInfo`)
}

/**
 * 检查 Codex 是否已安装
 */
export async function CheckCodexInstalled(): Promise<boolean> {
  return Call.ByName(`${SERVICE_NAME}.CheckCodexInstalled`)
}

/**
 * 检查依赖
 */
export async function CheckDependencies(): Promise<DependencyCheckResult> {
  return Call.ByName(`${SERVICE_NAME}.CheckDependencies`)
}

/**
 * 获取工具版本信息
 */
export async function GetToolVersions(): Promise<Record<string, string>> {
  return Call.ByName(`${SERVICE_NAME}.GetToolVersions`)
}

/**
 * 检查配置状态
 */
export async function CheckConfigurationStatus(): Promise<AutoConfigResult> {
  return Call.ByName(`${SERVICE_NAME}.CheckConfigurationStatus`)
}

/**
 * 获取默认的 Base URL
 */
export async function GetDefaultBaseURL(): Promise<string> {
  return Call.ByName(`${SERVICE_NAME}.GetDefaultBaseURL`)
}
