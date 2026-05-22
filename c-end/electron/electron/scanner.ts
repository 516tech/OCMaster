import { spawn, ChildProcess } from 'child_process'
import { app } from 'electron'
import { join } from 'path'
import { existsSync, readFileSync } from 'fs'

function getScannerPath(): string {
  if (app.isPackaged) {
    const platform = process.platform
    const ext = platform === 'win32' ? '.exe' : ''
    return join(process.resourcesPath, 'scanner', `ocmaster${ext}`)
  }
  const platform = process.platform
  const ext = platform === 'win32' ? '.exe' : ''
  return join(app.getAppPath(), '..', '..', 'bin', `ocmaster${ext}`)
}

export interface ScanResult {
  success: boolean
  data?: string
  error?: string
}

export function runScan(signal?: AbortSignal): Promise<ScanResult> {
  return new Promise((resolve) => {
    const scannerPath = getScannerPath()
    if (!existsSync(scannerPath)) {
      resolve({ success: false, error: `扫描器未找到: ${scannerPath}` })
      return
    }

    const proc: ChildProcess = spawn(scannerPath, ['scan', '--out', 'json'], {
      timeout: 15000,
      windowsHide: true
    })

    let stdout = ''
    let stderr = ''

    // Support AbortController cancellation
    if (signal) {
      const onAbort = () => {
        proc.kill('SIGTERM')
        resolve({ success: false, error: '扫描已取消' })
      }
      signal.addEventListener('abort', onAbort, { once: true })
      if (signal.aborted) {
        proc.kill('SIGTERM')
        resolve({ success: false, error: '扫描已取消' })
        return
      }
    }

    proc.stdout?.on('data', (chunk: Buffer) => {
      stdout += chunk.toString()
    })

    proc.stderr?.on('data', (chunk: Buffer) => {
      stderr += chunk.toString()
    })

    proc.on('close', (code) => {
      if (signal?.aborted) return
      if (code === 0 && stdout.trim()) {
        resolve({ success: true, data: stdout.trim() })
      } else {
        resolve({
          success: false,
          error: stderr || `扫描器退出码 ${code}`
        })
      }
    })

    proc.on('error', (err) => {
      if (signal?.aborted) return
      resolve({
        success: false,
        error: `无法启动扫描器: ${err.message}`
      })
    })
  })
}

export function getVersion(): string {
  try {
    const versionPath = join(app.getAppPath(), '..', '..', 'VERSION')
    return readFileSync(versionPath, 'utf-8').trim()
  } catch {
    return app.getVersion()
  }
}
