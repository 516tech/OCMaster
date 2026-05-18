import client from './client'

export const referenceApi = {
  getCpu: () => client.get('/reference/cpu'),
  getRam: () => client.get('/reference/ram'),
  getCooler: () => client.get('/reference/cooler'),
}
