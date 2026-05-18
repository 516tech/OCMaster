import client from './client'

export const hardwareApi = {
  getByCode: (code: string) => client.get(`/hardware/${code}`),
  delete: (code: string) => client.delete(`/hardware/${code}`),
}
