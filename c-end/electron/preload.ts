import { contextBridge, ipcRenderer } from 'electron'

contextBridge.exposeInMainWorld('ocmaster', {
  scanHardware: () => ipcRenderer.invoke('scan-hardware'),
  getConfig: () => ipcRenderer.invoke('get-config'),
  setConfig: (key: string, value: unknown) => ipcRenderer.invoke('set-config', key, value),
  uploadHardware: (data: unknown) => ipcRenderer.invoke('upload-hardware', data),
  deleteHardware: (code: string) => ipcRenderer.invoke('delete-hardware', code),
  exportTxt: (data: string) => ipcRenderer.invoke('export-txt', data),
})
