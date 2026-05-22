import { app } from 'electron'
import { readFileSync, writeFileSync, existsSync, mkdirSync } from 'fs'
import { join, dirname } from 'path'

export interface AppConfig {
  apiUrl: string
  language: string
  autoScan: boolean
}

const defaultConfig: AppConfig = {
  apiUrl: 'https://localhost/api/v1',
  language: 'zh-CN',
  autoScan: true
}

let configPath: string

export function initConfig(): void {
  const userData = app.getPath('userData')
  configPath = join(userData, 'config.json')
}

export function loadConfig(): AppConfig {
  try {
    if (existsSync(configPath)) {
      const raw = readFileSync(configPath, 'utf-8')
      return { ...defaultConfig, ...JSON.parse(raw) }
    }
  } catch {
    // ignore parse errors, use defaults
  }
  return { ...defaultConfig }
}

export function saveConfig(config: Partial<AppConfig>): void {
  const merged = { ...loadConfig(), ...config }
  try {
    const dir = dirname(configPath)
    if (!existsSync(dir)) mkdirSync(dir, { recursive: true })
    writeFileSync(configPath, JSON.stringify(merged, null, 2), 'utf-8')
  } catch {
    // fail silently — settings persist to userData
  }
}
