// Toast notification composable — adapted from ImHex Toast/Banner pattern
// Auto-expires after 4 seconds (matches ImHex toast duration)

import { ElNotification } from 'element-plus'

type ToastType = 'success' | 'error' | 'warning' | 'info'

const DEFAULT_DURATION = 4000

export function useToast() {
  function show(type: ToastType, message: string, title?: string) {
    const titles: Record<ToastType, string> = {
      success: '成功',
      error: '错误',
      warning: '警告',
      info: '提示'
    }
    ElNotification({
      type,
      title: title || titles[type],
      message,
      duration: DEFAULT_DURATION,
      position: 'bottom-right',
      customClass: type === 'error' ? 'toast-error' : ''
    })
  }

  return {
    success: (msg: string) => show('success', msg),
    error: (msg: string) => show('error', msg),
    warning: (msg: string) => show('warning', msg),
    info: (msg: string) => show('info', msg)
  }
}
