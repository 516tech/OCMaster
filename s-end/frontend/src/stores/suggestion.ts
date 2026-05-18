import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Suggestion } from '../types'
import { suggestionApi } from '../api/suggestion'

export const useSuggestionStore = defineStore('suggestion', () => {
  const current = ref<Suggestion | null>(null)
  const history = ref<Suggestion[]>([])
  const total = ref(0)

  async function create(data: Partial<Suggestion>) {
    await suggestionApi.create(data)
  }

  async function fetchById(id: number) {
    const res = await suggestionApi.getById(id)
    current.value = res.data
  }

  async function fetchHistory(offset = 0, limit = 20) {
    const res = await suggestionApi.getHistory(offset, limit)
    history.value = res.data.items
    total.value = res.data.total
  }

  async function downloadPdf(id: number) {
    const res = await suggestionApi.getPdf(id)
    const url = URL.createObjectURL(new Blob([res.data]))
    const a = document.createElement('a')
    a.href = url; a.download = 'report.pdf'; a.click()
    URL.revokeObjectURL(url)
  }

  return { current, history, total, create, fetchById, fetchHistory, downloadPdf }
})
