import axios from 'axios'

const client = axios.create({
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' }
})

export async function uploadHardware(
  baseUrl: string,
  hardwareInfo: Record<string, unknown>
): Promise<string | null> {
  const res = await client.post(`${baseUrl}/hardware/upload`, hardwareInfo)
  return res.data?.shareCode ?? null
}

export async function deleteHardware(
  baseUrl: string,
  code: string
): Promise<boolean> {
  await client.delete(`${baseUrl}/hardware/${code}`)
  return true
}
