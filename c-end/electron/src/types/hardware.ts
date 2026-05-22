export interface CpuInfo {
  model: string
  cores: number
  threads: number
  baseFreq: string
  biosVersion: string
}

export interface MotherboardInfo {
  brand: string
  model: string
  chipset: string
  biosVersion: string
}

export interface RamStick {
  capacity: string
  frequency: string
  timings: string
  dieType: string
}

export interface RamInfo {
  totalCapacity: string
  stickCount: number
  channelCount: number
  sticks: RamStick[]
}

export interface GpuInfo {
  model: string
  vram: string
}

export interface PsuInfo {
  ratedWattage: string
  source: string
}

export interface CoolerInfo {
  type: string
  source: string
}

export interface HardwareInfo {
  cpu: CpuInfo
  motherboard: MotherboardInfo
  ram: RamInfo
  gpu: GpuInfo
  psu: PsuInfo
  cooler: CoolerInfo
}
