export interface CpuInfo { model: string; cores: number; threads: number; base_freq: string; bios_version: string }
export interface MotherboardInfo { brand: string; model: string; chipset: string; bios_version: string }
export interface RamStick { capacity: string; frequency: string; timings: string; die_type: string }
export interface RamInfo { total_capacity: string; stick_count: number; sticks: RamStick[]; channel_count: number }
export interface GpuInfo { model: string; vram: string }
export interface PsuInfo { rated_wattage: string; source: string }
export interface CoolerInfo { type: string; source: string }
export interface HardwareInfo {
  cpu: CpuInfo; motherboard: MotherboardInfo; ram: RamInfo
  gpu: GpuInfo; psu: PsuInfo; cooler: CoolerInfo
}
export interface HardwareUpload {
  id: number; share_code: string; hardware_data: HardwareInfo; expires_at: string; created_at: string
}
export interface Merchant {
  id: number; phone: string; status: string; risk_template: string
}
export interface CpuSuggestion { frequency: string; voltage: string; notes: string }
export interface RamSuggestion { frequency: string; timings: string; voltage: string; notes: string }
export interface StabilityTest { tool: string; duration: string }
export interface Suggestion {
  id: number; merchant_id: number; hardware_id: number
  cpu_suggestion: CpuSuggestion; ram_suggestion: RamSuggestion
  stability_test: StabilityTest; risk_warning: string
}
