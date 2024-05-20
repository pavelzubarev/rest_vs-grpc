using Google.Protobuf.Collections;
using Grpc.Core;
using GrpcServer;

namespace grpc.Services;

public class GrpcPolygonService : PolygonService.PolygonServiceBase
{
    public override Task<PolygonResponse> CalculateArea(PolygonRequest request, ServerCallContext context)
    {
        var area = CalculatePolygonArea(request.Points);
        return Task.FromResult(new PolygonResponse { Area = area });
    }

    private static double CalculatePolygonArea(RepeatedField<Point> points)
    {
        var n = points.Count;
        double area = 0;

        for (var i = 0; i < n; i++)
        {
            var current = points[i];
            var next = points[(i + 1) % n];
            area += current.X * next.Y - next.X * current.Y;
        }

        return 0.5 * Math.Abs(area);
    }
}