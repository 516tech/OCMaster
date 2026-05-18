using System.Text.Json;

namespace OCMaster.App.Services;

public class AppConfig
{
    public string ApiUrl { get; set; } = "https://localhost/api/v1";
    public string Language { get; set; } = "zh-CN";
    public bool AutoScan { get; set; } = true;
}

public static class ConfigService
{
    private static readonly string ConfigPath =
        Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
                     "OCMaster", "config.json");

    private static readonly JsonSerializerOptions JsonOpts = new()
    {
        WriteIndented = true,
        PropertyNameCaseInsensitive = true
    };

    private static AppConfig? _config;

    public static AppConfig Load()
    {
        try
        {
            if (File.Exists(ConfigPath))
            {
                var json = File.ReadAllText(ConfigPath);
                _config = JsonSerializer.Deserialize<AppConfig>(json, JsonOpts) ?? new AppConfig();
            }
        }
        catch { /* 使用默认配置 */ }

        _config ??= new AppConfig();
        return _config;
    }

    public static void Save(AppConfig config)
    {
        _config = config;
        var dir = Path.GetDirectoryName(ConfigPath);
        if (!string.IsNullOrEmpty(dir)) Directory.CreateDirectory(dir);
        File.WriteAllText(ConfigPath, JsonSerializer.Serialize(config, JsonOpts));
    }

    public static void Set(string key, object value)
    {
        _config ??= Load();
        var prop = typeof(AppConfig).GetProperty(key);
        if (prop != null)
        {
            prop.SetValue(_config, Convert.ChangeType(value, prop.PropertyType));
            Save(_config);
        }
    }

    public static object? Get(string key)
    {
        _config ??= Load();
        return typeof(AppConfig).GetProperty(key)?.GetValue(_config);
    }
}
