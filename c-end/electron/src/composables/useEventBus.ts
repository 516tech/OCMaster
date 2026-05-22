// Typed pub/sub event bus — adapted from ImHex EventManager pattern
// Events: announce state changes (frameBegin, scanComplete, themeChanged)
// Requests: ask for actions (openView, showToast, cancelScan)

type Handler<T = void> = (payload: T) => void
type Unsubscribe = () => void

class EventChannel<T = void> {
  private handlers = new Set<Handler<T>>()

  subscribe(fn: Handler<T>): Unsubscribe {
    this.handlers.add(fn)
    return () => this.handlers.delete(fn)
  }

  post(payload: T): void {
    for (const fn of this.handlers) fn(payload)
  }
}

// --- Application Events ---
export const events = {
  scanStarted: new EventChannel<void>(),
  scanCompleted: new EventChannel<{ success: boolean; error?: string }>(),
  themeChanged: new EventChannel<'dark' | 'light'>(),
  configChanged: new EventChannel<Record<string, unknown>>(),
  viewOpened: new EventChannel<string>(),
  viewClosed: new EventChannel<string>()
}

// --- Application Requests ---
export const requests = {
  showToast: new EventChannel<{ type: 'success' | 'error' | 'warning' | 'info'; message: string }>(),
  cancelScan: new EventChannel<void>(),
  navigateTo: new EventChannel<string>()
}

export function useEventBus() {
  return { events, requests }
}
