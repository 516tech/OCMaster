import { app, BrowserWindow, shell } from 'electron'
import { join } from 'path'
import { registerIpcHandlers } from './ipc-handlers'

let mainWindow: BrowserWindow | null = null

function createWindow(): void {
  const saved = loadWindowState()

  mainWindow = new BrowserWindow({
    x: saved.x,
    y: saved.y,
    width: saved.width || 960,
    height: saved.height || 680,
    minWidth: 800,
    minHeight: 600,
    title: '超频大师 OCMaster',
    webPreferences: {
      preload: join(__dirname, '../preload/preload.js'),
      contextIsolation: true,
      nodeIntegration: false,
      sandbox: false
    }
  })

  if (saved.maximized) mainWindow.maximize()

  // Save window state on move/resize
  mainWindow.on('close', () => saveWindowState(mainWindow!))

  mainWindow.on('ready-to-show', () => {
    mainWindow?.show()
  })

  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    shell.openExternal(url)
    return { action: 'deny' }
  })

  if (process.env.NODE_ENV === 'development') {
    mainWindow.loadURL('http://localhost:5173')
    mainWindow.webContents.openDevTools()
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

// --- Window state persistence (adapted from ImHex LayoutManager) ---
import { readFileSync, writeFileSync, existsSync, mkdirSync } from 'fs'

interface WindowState {
  x?: number; y?: number
  width: number; height: number
  maximized: boolean
}

function statePath(): string {
  const dir = join(app.getPath('userData'))
  return join(dir, 'window-state.json')
}

function loadWindowState(): WindowState {
  try {
    const raw = readFileSync(statePath(), 'utf-8')
    return JSON.parse(raw)
  } catch {
    return { width: 960, height: 680, maximized: false }
  }
}

function saveWindowState(win: BrowserWindow): void {
  const maximized = win.isMaximized()
  const bounds = maximized ? { x: 0, y: 0, width: 0, height: 0 } : win.getBounds()
  const dir = join(app.getPath('userData'))
  try {
    if (!existsSync(dir)) mkdirSync(dir, { recursive: true })
    writeFileSync(statePath(), JSON.stringify({
      x: bounds.x, y: bounds.y,
      width: bounds.width || 960, height: bounds.height || 680,
      maximized
    }))
  } catch { /* ignore */ }
}

const gotTheLock = app.requestSingleInstanceLock()
if (!gotTheLock) {
  app.quit()
} else {
  app.on('second-instance', () => {
    if (mainWindow) {
      if (mainWindow.isMinimized()) mainWindow.restore()
      mainWindow.focus()
    }
  })
}

app.whenReady().then(() => {
  registerIpcHandlers()
  createWindow()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit()
})
