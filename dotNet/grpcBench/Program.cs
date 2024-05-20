using System.Diagnostics;
using Grpc.Net.Client;
using GrpcServer;
using Microsoft.Extensions.Logging;

namespace grpcBench;

public static class Program
{
    private const int NumberOfClients = 100;
    private const int RequestsPerClient = 1000;

    private static int _successCount;
    private static int _failureCount;
    private static readonly object _lock = new();

    public static async Task Main(string[] args)
    {
        var loggerFactory = LoggerFactory.Create(builder =>
        {
            builder
                .AddFilter("Grpc", LogLevel.Debug)
                .AddConsole();
        });

        var logger = loggerFactory.CreateLogger("benchmark");

        var points = new[]
        {
            new Point { X = 0, Y = 0 },
            new Point { X = 4, Y = 0 },
            new Point { X = 4, Y = 3 }
        };

        var request = new PolygonRequest();
        request.Points.AddRange(points);

        var stopwatch = Stopwatch.StartNew();

        var tasks = Enumerable.Range(0, NumberOfClients)
            .Select(_ => Task.Run(() => RunClientRequests(request, logger)))
            .ToArray();

        await Task.WhenAll(tasks);

        stopwatch.Stop();

        var totalRequests = NumberOfClients * RequestsPerClient;
        var duration = stopwatch.Elapsed;
        var requestsPerSecond = totalRequests / duration.TotalSeconds;

        Console.WriteLine($"Test completed in {duration}");
        Console.WriteLine($"Total requests: {totalRequests}");
        Console.WriteLine($"Successful requests: {_successCount}");
        Console.WriteLine($"Failed requests: {_failureCount}");
        Console.WriteLine($"Requests per second: {requestsPerSecond}");
    }

    private static async Task RunClientRequests(PolygonRequest request, ILogger logger)
    {
        using var channel = GrpcChannel.ForAddress("http://localhost:5001");
        var client = new PolygonService.PolygonServiceClient(channel);

        for (var i = 0; i < RequestsPerClient; i++)
            try
            {
                await RetryPolicy(async () => { await client.CalculateAreaAsync(request); }, logger);

                lock (_lock)
                {
                    _successCount++;
                }
            }
            catch (Exception ex)
            {
                lock (_lock)
                {
                    _failureCount++;
                }

                logger.LogError("Request {I}: An error occurred: {ExMessage}", i + 1, ex.Message);
            }
    }

    private static async Task RetryPolicy(Func<Task> action, ILogger logger, int maxRetryAttempts = 3,
        int delayMilliseconds = 200)
    {
        var retryAttempts = 0;
        while (true)
            try
            {
                await action();
                return;
            }
            catch (Exception ex) when (retryAttempts < maxRetryAttempts)
            {
                retryAttempts++;
                logger.LogWarning("Attempt {RetryAttempts} failed: {ExMessage}. Retrying in {DelayMilliseconds}ms...",
                    retryAttempts, ex.Message, delayMilliseconds);
                await Task.Delay(delayMilliseconds);
            }
    }
}