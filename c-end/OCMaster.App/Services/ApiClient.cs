using System.Net.Http.Json;
using OCMaster.App.Models;

namespace OCMaster.App.Services;

public class ApiClient
{
    private readonly HttpClient _http;
    private string _baseUrl;

    public ApiClient(string baseUrl)
    {
        _baseUrl = baseUrl.TrimEnd('/');
        _http = new HttpClient { Timeout = TimeSpan.FromSeconds(10) };
    }

    public async Task<string?> UploadHardware(HardwareInfo info)
    {
        var resp = await _http.PostAsJsonAsync($"{_baseUrl}/hardware/upload", info);
        if (!resp.IsSuccessStatusCode) return null;
        var result = await resp.Content.ReadFromJsonAsync<UploadResponse>();
        return result?.ShareCode;
    }

    public async Task<bool> DeleteHardware(string code)
    {
        var resp = await _http.DeleteAsync($"{_baseUrl}/hardware/{code}");
        return resp.IsSuccessStatusCode;
    }

    private class UploadResponse
    {
        public string? ShareCode { get; set; }
    }
}
