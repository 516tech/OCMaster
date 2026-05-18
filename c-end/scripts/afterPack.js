// electron-builder afterPack hook - 裁剪 Chromium 冗余文件
const fs = require('fs')
const path = require('path')

exports.default = async function (context) {
  const appDir = context.appOutDir
  const removeFiles = [
    'vk_swiftshader.dll', 'vk_swiftshader_icd.json', 'vulkan-1.dll',     // Vulkan 软件渲染 ~6MB
    'd3dcompiler_47.dll',                                                  // Direct3D 编译器 ~4.7MB
  ]
  let removed = 0
  for (const file of removeFiles) {
    const filePath = path.join(appDir, file)
    if (fs.existsSync(filePath)) {
      const size = fs.statSync(filePath).size
      fs.unlinkSync(filePath)
      removed += size
      console.log(`Removed: ${file} (${(size / 1024 / 1024).toFixed(1)}MB)`)
    }
  }
  console.log(`Total removed: ${(removed / 1024 / 1024).toFixed(1)}MB`)
}
