import { app, BrowserWindow, ipcMain, Menu } from 'electron'
import path from 'path'
import Store from 'electron-store'

// ============ 反调试检测 ============
function antiDebug() {
  // 检测 --inspect / --inspect-brk 启动参数
  for (const arg of process.argv) {
    if (arg.startsWith('--inspect')) {
      app.quit()
      return
    }
  }
  // 检测 ELECTRON_RUN_AS_NODE 环境变量
  if (process.env.ELECTRON_RUN_AS_NODE) {
    app.quit()
    return
  }
}

antiDebug()

const store = new Store({
  defaults: {
    apiUrl: 'https://localhost/api/v1',
    language: 'zh-CN',
    autoScan: true,
  },
})

function createWindow() {
  const win = new BrowserWindow({
    width: 960,
    height: 680,
    title: '超频大师 OCMaster',
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false,
      devTools: false, // 禁用 DevTools
    },
  })

  // 阻止打开 DevTools 的快捷键
  win.webContents.on('before-input-event', (_e, input: any) => {
    const blocked = ((input.control || input.meta) && ['i','j','u'].includes(input.key)) || input.key === 'F12'
    if (blocked && (input as any).preventDefault) (input as any).preventDefault()
  })

  if (process.env.VITE_DEV_SERVER_URL) {
    win.loadURL(process.env.VITE_DEV_SERVER_URL)
  } else {
    win.loadFile(path.join(__dirname, '../dist/index.html'))
  }
}

// Hardware scan: native addon 打包在 app 内
// 开发时路径: ../../build/Release/, 打包后路径: ../../native/
function scanHardware() {
  const paths = [
    path.join(__dirname, '../../native/hardware_scanner.node'),                     // 打包后
    path.join(__dirname, '../electron/native/build/Release/hardware_scanner.node'),  // 开发时
  ]
  for (const p of paths) {
    try { return require(p).scan() } catch { /* try next path */ }
  }
  throw new Error('扫描模块不可用')
}

// IPC handlers
ipcMain.handle('scan-hardware', async () => scanHardware())
ipcMain.handle('get-config', async () => store.store)
ipcMain.handle('set-config', async (_e, key: string, value: unknown) => { store.set(key, value); return true })
ipcMain.handle('upload-hardware', async (_e, hardwareData: unknown) => {
  const apiUrl = store.get('apiUrl')
  const resp = await fetch(`${apiUrl}/hardware/upload`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(hardwareData),
  })
  return resp.json()
})
ipcMain.handle('delete-hardware', async (_e, code: string) => {
  const apiUrl = store.get('apiUrl')
  const resp = await fetch(`${apiUrl}/hardware/${code}`, { method: 'DELETE' })
  return resp.json()
})

app.whenReady().then(() => {
  Menu.setApplicationMenu(null) // 移除英文默认菜单
  createWindow()
  app.on('activate', () => { if (BrowserWindow.getAllWindows().length === 0) createWindow() })
})
app.on('window-all-closed', () => { if (process.platform !== 'darwin') app.quit() })
