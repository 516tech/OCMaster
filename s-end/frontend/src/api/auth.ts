import client from './client'

export const authApi = {
  login: (phone: string, password: string) => client.post('/merchant/login', { phone, password }),
  register: (phone: string, password: string) => client.post('/merchant/register', { phone, password }),
}

export const profileApi = {
  get: () => client.get('/merchant/profile'),
  update: (data: any) => client.put('/merchant/profile', data),
  changePassword: (old_password: string, new_password: string) => client.put('/merchant/password', { old_password, new_password }),
}
