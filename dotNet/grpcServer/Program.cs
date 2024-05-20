using grpc.Services;
using Microsoft.AspNetCore.Server.Kestrel.Core;

var builder = WebApplication.CreateBuilder(args);

// Add gRPC services to the container
builder.Services.AddGrpc();

// Configure Kestrel to listen on port 5001
builder.WebHost.ConfigureKestrel(options => { options.ListenLocalhost(5001, o => o.Protocols = HttpProtocols.Http2); });

var app = builder.Build();

app.MapGrpcService<GrpcPolygonService>();

// Optional: Configure HTTP request pipeline if needed
app.MapGet("/",
    async context =>
    {
        await context.Response.WriteAsync("Communication with gRPC endpoints must be made through a gRPC client.");
    });

app.Run();