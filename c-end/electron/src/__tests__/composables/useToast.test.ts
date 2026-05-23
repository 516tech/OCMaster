import { describe, it, expect, vi } from 'vitest'
import { useToast } from '@/composables/useToast'

// Track ElNotification calls
const notifCalls: Array<Record<string, unknown>> = []
vi.mock('element-plus', () => ({
  ElNotification: (opts: Record<string, unknown>) => { notifCalls.push(opts) }
}))

describe('useToast', () => {
  it('success toast calls ElNotification with success type', () => {
    notifCalls.length = 0
    const toast = useToast()
    toast.success('扫描完成')
    expect(notifCalls.length).toBe(1)
    expect(notifCalls[0].type).toBe('success')
    expect(notifCalls[0].message).toBe('扫描完成')
  })

  it('error toast calls ElNotification with error type', () => {
    notifCalls.length = 0
    const toast = useToast()
    toast.error('扫描失败')
    expect(notifCalls.length).toBe(1)
    expect(notifCalls[0].type).toBe('error')
    expect(notifCalls[0].message).toBe('扫描失败')
  })

  it('warning toast calls ElNotification with warning type', () => {
    notifCalls.length = 0
    const toast = useToast()
    toast.warning('已取消')
    expect(notifCalls[0].type).toBe('warning')
    expect(notifCalls[0].message).toBe('已取消')
  })

  it('info toast calls ElNotification with info type', () => {
    notifCalls.length = 0
    const toast = useToast()
    toast.info('提示信息')
    expect(notifCalls[0].type).toBe('info')
    expect(notifCalls[0].message).toBe('提示信息')
  })

  it('all toasts set duration to 4000 and position bottom-right', () => {
    notifCalls.length = 0
    const toast = useToast()
    toast.success('ok')
    expect(notifCalls[0].duration).toBe(4000)
    expect(notifCalls[0].position).toBe('bottom-right')
  })
})
