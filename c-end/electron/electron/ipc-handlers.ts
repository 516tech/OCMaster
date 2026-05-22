import { ipcMain, dialog, BrowserWindow } from 'electron'
import { runScan } from './scanner'
import { loadConfig, saveConfig, initConfig, AppConfig } from './config'
import { readFileSync, writeFileSync, existsSync, mkdirSync } from 'fs'
import { join } from 'path'
import { app } from 'electron'

function lastTabPath(): string {
  const dir = join(app.getPath('userData'))
  return join(dir, 'last-tab.json')
}

export function registerIpcHandlers(): void {
  initConfig()

  ipcMain.handle('scan:all', async () => {
    return runScan()
  })

  ipcMain.handle('config:load', () => {
    return loadConfig()
  })

  ipcMain.handle('config:save', (_event, config: Partial<AppConfig>) => {
    saveConfig(config)
    return { ok: true }
  })

  ipcMain.handle('export:txt', async (_event, content: string) => {
    const win = BrowserWindow.getFocusedWindow()
    if (!win) return { success: false, error: 'No window' }

    const result = await dialog.showSaveDialog(win, {
      defaultPath: `ocmaster-hardware-${new Date().toISOString().slice(0, 10)}.txt`,
      filters: [{ name: 'Text', extensions: ['txt'] }]
    })

    if (result.canceled || !result.filePath) {
      return { success: false, canceled: true }
    }

    try {
      writeFileSync(result.filePath, content, 'utf-8')
      return { success: true, path: result.filePath }
    } catch (err) {
      return { success: false, error: String(err) }
    }
  })

  ipcMain.handle('app:version', () => {
    try {
      return readFileSync(join(app.getAppPath(), '..', '..', 'VERSION'), 'utf-8').trim()
    } catch {
      return app.getVersion()
    }
  })

  // --- Window state / last-tab persistence (ImHex LayoutManager pattern) ---
  ipcMain.handle('window:get-last-tab', () => {
    try {
      return JSON.parse(readFileSync(lastTabPath(), 'utf-8')).tab || '/'
    } catch {
      return '/'
    }
  })

  ipcMain.handle('window:save-last-tab', (_event, tab: string) => {
    try {
      const dir = join(app.getPath('userData'))
      if (!existsSync(dir)) mkdirSync(dir, { recursive: true })
      writeFileSync(lastTabPath(), JSON.stringify({ tab }))
    } catch { /* ignore */ }
    return { ok: true }
  })
}
