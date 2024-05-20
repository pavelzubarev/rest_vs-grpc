using System.Net.Http.Json;
using RestServer.Models;

namespace restClient;

public static class Program
{
    private static async Task Main(string[] args)
    {
        var points = new List<Point>
        {
            new() { X = 0, Y = 0 },
            new() { X = 4, Y = 0 },
            new() { X = 4, Y = 3 }
        };

        var polygonRequest = new PolygonRequest { Points = points };

        using var httpClient = new HttpClient();
        httpClient.BaseAddress = new Uri("http://localhost:5002");

        try
        {
            var response = await httpClient.PostAsJsonAsync("/api/polygon/calculate-area", polygonRequest);

            if (response.IsSuccessStatusCode)
            {
                var polygonResponse = await response.Content.ReadFromJsonAsync<PolygonResponse>();
                Console.WriteLine($"The area of the polygon is: {polygonResponse?.Area}");
            }
            else
            {
                Console.WriteLine($"Error: {response.StatusCode}");
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Exception: {ex.Message}");
        }
    }
}