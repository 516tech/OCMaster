export interface ScanResult {
  success: boolean
  data?: string
  error?: string
}

export interface ExportResult {
  success: boolean
  path?: string
  error?: string
  canceled?: boolean
}

export interface ElectronAPI {
  scanAll: () => Promise<ScanResult>
  getConfig: () => Promise<Record<string, unknown>>
  saveConfig: (config: Record<string, unknown>) => Promise<{ ok: boolean }>
  exportTxt: (content: string) => Promise<ExportResult>
  getVersion: () => Promise<string>
  getScannerPath: () => Promise<string>
  getLastTab: () => Promise<string>
  saveLastTab: (tab: string) => Promise<{ ok: boolean }>
  getWindowState: () => Promise<Record<string, unknown>>
  saveWindowState: (state: Record<string, unknown>) => Promise<{ ok: boolean }>
}

declare global {
  interface Window {
    electronAPI: ElectronAPI
  }
}

export {}
