import { contextBridge, ipcRenderer } from 'electron'

contextBridge.exposeInMainWorld('electronAPI', {
  scanAll: () => ipcRenderer.invoke('scan:all'),
  getConfig: () => ipcRenderer.invoke('config:load'),
  saveConfig: (config: Record<string, unknown>) => ipcRenderer.invoke('config:save', config),
  exportTxt: (content: string) => ipcRenderer.invoke('export:txt', content),
  getVersion: () => ipcRenderer.invoke('app:version'),
  getScannerPath: () => ipcRenderer.invoke('scanner:path'),

  // Window state persistence (ImHex LayoutManager pattern)
  getLastTab: () => ipcRenderer.invoke('window:get-last-tab'),
  saveLastTab: (tab: string) => ipcRenderer.invoke('window:save-last-tab', tab),
  getWindowState: () => ipcRenderer.invoke('window:get-state'),
  saveWindowState: (state: Record<string, unknown>) => ipcRenderer.invoke('window:save-state', state)
})
