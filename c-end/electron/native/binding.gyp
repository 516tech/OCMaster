{
  "targets": [
    {
      "target_name": "hardware_scanner",
      "sources": [
        "src/binding.cpp"
      ],
      "conditions": [
        ["OS=='win'", { "sources": ["src/hardware_win.cpp"], "libraries": ["ole32.lib", "oleaut32.lib"] }],
        ["OS=='linux'", { "sources": ["src/hardware_linux.cpp"] }],
        ["OS=='mac'", { "sources": ["src/hardware_mac.cpp"], "libraries": ["-framework IOKit", "-framework Foundation"] }]
      ],
      "include_dirs": ["<!@(node -p \"require('node-addon-api').include\")"],
      "defines": ["NAPI_DISABLE_CPP_EXCEPTIONS"]
    }
  ]
}
