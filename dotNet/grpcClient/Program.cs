using Grpc.Net.Client;
using GrpcServer;

namespace grpcClient;

public static class Program
{
    private static async Task Main(string[] args)
    {
        // The port number must match the port of the gRPC server.
        using var channel = GrpcChannel.ForAddress("http://localhost:5001");
        var client = new PolygonService.PolygonServiceClient(channel);

        var request = new PolygonRequest
        {
            Points =
            {
                new Point { X = 0, Y = 0 },
                new Point { X = 4, Y = 0 },
                new Point { X = 4, Y = 3 }
            }
        };

        var response = await client.CalculateAreaAsync(request);

        Console.WriteLine($"The area of the polygon is: {response.Area}");
    }
}