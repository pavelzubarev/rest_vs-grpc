using System.Diagnostics;
using System.Net.Http.Json;
using RestServer.Models;

namespace restBench;

public static class Program
{
    private static async Task Main(string[] args)
    {
        const int clientCount = 100;
        const int requestsPerClient = 1000;
        const string apiUrl = "http://localhost:5002/api/polygon/calculate-area";
        var points = new List<Point>
        {
            new() { X = 0, Y = 0 },
            new() { X = 4, Y = 0 },
            new() { X = 4, Y = 3 },
        };

        var polygonRequest = new PolygonRequest { Points = points };

        var totalRequests = clientCount * requestsPerClient;
        var successfulRequests = 0;
        var failedRequests = 0;

        var stopwatch = Stopwatch.StartNew();

        var tasks = new List<Task>();

        for (var i = 0; i < clientCount; i++)
            tasks.Add(Task.Run(async () =>
            {
                using var httpClient = new HttpClient();
                for (var j = 0; j < requestsPerClient; j++)
                    try
                    {
                        var response = await httpClient.PostAsJsonAsync(apiUrl, polygonRequest);
                        if (response.IsSuccessStatusCode)
                        {
                            var polygonResponse = await response.Content.ReadFromJsonAsync<PolygonResponse>();
                            if (polygonResponse != null)
                                Interlocked.Increment(ref successfulRequests);
                            else
                                Interlocked.Increment(ref failedRequests);
                        }
                        else
                        {
                            Interlocked.Increment(ref failedRequests);
                        }
                    }
                    catch
                    {
                        Interlocked.Increment(ref failedRequests);
                    }
            }));

        await Task.WhenAll(tasks);

        stopwatch.Stop();

        var successRate = (double)successfulRequests / totalRequests * 100;
        var failureRate = (double)failedRequests / totalRequests * 100;
        var requestsPerSecond = totalRequests / stopwatch.Elapsed.TotalSeconds;

        Console.WriteLine($"Test completed in {stopwatch.Elapsed}");
        Console.WriteLine($"Total requests: {totalRequests}");
        Console.WriteLine($"Successful requests: {successfulRequests}");
        Console.WriteLine($"Failed requests: {failedRequests}");
        Console.WriteLine($"Success rate: {successRate:F2}%");
        Console.WriteLine($"Failure rate: {failureRate:F2}%");
        Console.WriteLine($"Requests per second: {requestsPerSecond:F6}");
    }
}