import { describe, it, expect, vi } from 'vitest'
import { events, requests, useEventBus } from '@/composables/useEventBus'

describe('useEventBus', () => {
  it('returns events and requests', () => {
    const bus = useEventBus()
    expect(bus.events).toBeDefined()
    expect(bus.requests).toBeDefined()
  })
})

describe('EventChannel', () => {
  it('subscribe and post', () => {
    const fn = vi.fn()
    events.scanStarted.subscribe(fn)
    events.scanStarted.post()
    expect(fn).toHaveBeenCalledOnce()
  })

  it('post with payload', () => {
    const fn = vi.fn()
    events.scanCompleted.subscribe(fn)
    events.scanCompleted.post({ success: true, error: undefined })
    expect(fn).toHaveBeenCalledWith({ success: true, error: undefined })
  })

  it('unsubscribe stops delivery', () => {
    const fn = vi.fn()
    const unsub = events.scanStarted.subscribe(fn)
    events.scanStarted.post()
    expect(fn).toHaveBeenCalledOnce()
    unsub()
    events.scanStarted.post()
    expect(fn).toHaveBeenCalledOnce() // still 1
  })

  it('multiple handlers all receive event', () => {
    const a = vi.fn()
    const b = vi.fn()
    events.scanStarted.subscribe(a)
    events.scanStarted.subscribe(b)
    events.scanStarted.post()
    expect(a).toHaveBeenCalledOnce()
    expect(b).toHaveBeenCalledOnce()
  })

  it('Requests channel works independently', () => {
    const fn = vi.fn()
    requests.cancelScan.subscribe(fn)
    requests.cancelScan.post()
    expect(fn).toHaveBeenCalledOnce()
  })

  it('themeChanged event delivers dark/light', () => {
    const fn = vi.fn()
    events.themeChanged.subscribe(fn)
    events.themeChanged.post('dark')
    expect(fn).toHaveBeenCalledWith('dark')
    events.themeChanged.post('light')
    expect(fn).toHaveBeenCalledWith('light')
  })
})
