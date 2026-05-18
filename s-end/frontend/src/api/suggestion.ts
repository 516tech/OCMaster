import client from './client'

export const suggestionApi = {
  create: (data: any) => client.post('/suggestions', data),
  getById: (id: number) => client.get(`/suggestions/${id}`),
  getHistory: (offset: number, limit: number) => client.get('/suggestions/history', { params: { offset, limit } }),
  getPdf: (id: number) => client.get(`/suggestions/${id}/pdf`, { responseType: 'blob' }),
}
