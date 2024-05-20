using Microsoft.AspNetCore.Mvc;
using RestServer.Models;

namespace RestServer.Conrtollers;

[ApiController]
[Route("api/[controller]")]
public class PolygonController : ControllerBase
{
    [HttpPost("calculate-area")]
    public ActionResult<PolygonResponse> CalculateArea([FromBody] PolygonRequest polygonRequest)
    {
        if (polygonRequest.Points.Count < 3) return BadRequest("A polygon must have at least 3 points.");

        var area = CalculatePolygonArea(polygonRequest.Points);
        return Ok(new PolygonResponse { Area = area });
    }

    private double CalculatePolygonArea(List<Point> points)
    {
        var n = points.Count;
        double area = 0;

        for (var i = 0; i < n; i++)
        {
            var j = (i + 1) % n;
            area += points[i].X * points[j].Y;
            area -= points[j].X * points[i].Y;
        }

        area = Math.Abs(area) / 2.0;
        return area;
    }
}