文档版本：V4.1
更新日期：2026-05-22
文档状态：设计稿
所属项目：超频大师（OCMaster）— UI 增强 (ImHex 模式)

一、评估维度

| 维度 | 通过标准 | 当前状态 |
|------|---------|---------|
| 功能完整性 | PRD checkbox 全部勾选 | 阶段〇 8/10 |
| 架构合规 | 符合 architecture.md V4.1 | EventBus + Theme + Toast |
| Go 测试 | scanner 80%+ | 80.7% ✅ |
| Electron 编译 | npm run build 通过 | ✅ (main 6.7KB + preload 0.9KB + renderer) |
| CI Windows | Go CLI + JSON test + electron-builder | ✅ (单 job, windows-latest) |
| Go 命名规范 | ocmaster_${VERSION}_windows_amd64.exe | ✅ |

二、CI 修复 (本轮)

| 问题 | 修复 |
|------|------|
| Matrix 三平台 → Windows only | 简化为单 `build` job |
| Go CLI 未带版本号 | Rename → `ocmaster_0.0.1_windows_amd64.exe` |
| `extraResources: ../bin/ocmaster` 无法匹配 `.exe` | 修正为 `../bin/ocmaster.exe` → `scanner/ocmaster.exe` |
| Release 下载所有平台产物 | 仅下载 `OCMaster-Windows` |
| package.json 版本与 VERSION 不一致 | CI 中 `node -e` 同步版本再构建 |

二、本轮新增 (ImHex UI 模式)

| 模块 | 参考 ImHex | 文件 | 状态 |
|------|----------|------|------|
| Typed EventBus (Event/Request) | EventManager | src/composables/useEventBus.ts | ✅ |
| 暗色/亮色主题切换 | ThemeManager | src/composables/useTheme.ts + stores/theme.ts | ✅ |
| Toast 通知 (4s 自动消失) | Toast/Banner | src/composables/useToast.ts | ✅ |
| 扫描进度 + 取消 (AbortController) | TaskManager | src/stores/scan.ts | ✅ |
| 窗口状态记忆 (位置/大小) | LayoutManager | electron/main.ts window-state.json | ✅ |
| 最后标签页记忆 | LayoutManager | electron/ipc-handlers.ts last-tab.json | ✅ |
| 全局错误 → Toast | EventManager | src/main.ts errorHandler | ✅ |

三、构建产物对比

| | 旧 (V4.0) | 新 (V4.1) |
|---|---|---|
| Main process | 5.07 KB | 6.83 KB |
| Preload | 0.52 KB | 0.90 KB |
| Renderer chunks | 6 files | 9 files (新增 useToast/SettingsPage+) |
| 编译结果 | ✅ PASS | ✅ PASS
