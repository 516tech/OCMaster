export const defaultConfig = {
  apiUrl: 'https://localhost/api/v1',
  language: 'zh-CN',
  autoScan: true,
} as const

export type AppConfig = typeof defaultConfig
