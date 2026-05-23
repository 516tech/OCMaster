<h1 align="center">OCMaster Evaluate</h1>

<p align="center">
    本轮迭代: prd.md + architecture.md 重写 (ImHex 风格)
</p>

## Goals

| # | 目标 | 标准 | 结果 |
|---|------|------|------|
| 1 | prd.md 重写 | ImHex 风格, ≤200行, ≤3表, ≤5同级标题 | ✅ 172行, 4H2, 3表 |
| 2 | architecture.md 重写 | 与 prd.md 同风格, 代码一致 | ✅ 151行, 5H2, 2表 |
| 3 | 文件树一致性 | 列出的 25 个文件全部存在 | ✅ 25/25 |
| 4 | IPC 协议一致性 | 文档 channel 与代码匹配 | ✅ 7/7 |
| 5 | CI 流程一致性 | 文档流程与 build-c-end.yml 匹配 | ✅ |

## Tests

| 测试 | 命令 | 结果 |
|------|------|------|
| Go 单元测试 | `go test ./scanner/... -cover` | 80.7% PASS |
| Vue 单元测试 | `npm test` (vitest) | 25/25 PASS (5 suites) |
| Electron 编译 | `npx electron-vite build` | PASS (6.7K+0.9K+9 chunks) |

### Test Coverage Detail

| Suite | Tests | Coverage |
|-------|-------|----------|
| stores/scan.test.ts | 7 | scan, cancel, reset, progress, events, error |
| stores/settings.test.ts | 3 | load defaults, load partial, save |
| composables/useEventBus.test.ts | 7 | subscribe, post, unsubscribe, payload, Event vs Request |
| composables/useTheme.test.ts | 3 | load preference, default, persist |
| composables/useToast.test.ts | 5 | success/error/warning/info + duration/position |

## Bugs

- [x] architecture.md 7 个 H2 → 合并为 5 个
- [x] architecture.md 与 prd.md 内容重叠 → 改为互补
- [x] IPC channel 文档与实际代码核对 → 全部匹配
- [ ] 三平台 GUI 手动测试 (需真实硬件)
- [ ] macOS 代码签名 (需 Apple 开发者账号)

## Summary

本轮完成 prd.md + architecture.md 的 ImHex 风格重写:
- 统一视觉语言 (居中标题, ASCII 架构图, 简洁分层)
- architecture.md 聚焦技术实现细节 (File Tree / IPC / Data Flow / CI / Decisions)
- prd.md 聚焦产品特性描述 (Features / Getting Started)
- 两份文档互补, 无重叠内容
