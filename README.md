## Basic REST vs gRPC comparison

* This code is an attempt to compare plain REST vs gRPC services performance.
* Each service calculates polygon surface
* To saturate each server implementation: 100 parallel clients, 1000 sequential requests per client


Results are taken on Apple M2 pro chip. 
```shell

% go test -v

=== RUN   TestComparePerformance
=== RUN   TestComparePerformance/gRPC_Performance
2024/05/20 12:51:21 gRPC Test completed in 1.033512875s
2024/05/20 12:51:21 gRPC Total requests: 100000
2024/05/20 12:51:21 gRPC Successful requests: 100000
2024/05/20 12:51:21 gRPC Failed requests: 0
2024/05/20 12:51:21 gRPC Success rate: 100.00%
2024/05/20 12:51:21 gRPC Failure rate: 0.00%
2024/05/20 12:51:21 gRPC Requests per second: 96757.381953
=== RUN   TestComparePerformance/REST_Performance
2024/05/20 12:52:41 REST Test completed in 1m19.474647708s
2024/05/20 12:52:41 REST Total requests: 100000
2024/05/20 12:52:41 REST Successful requests: 79544
2024/05/20 12:52:41 REST Failed requests: 20456
2024/05/20 12:52:41 REST Success rate: 79.54%
2024/05/20 12:52:41 REST Failure rate: 20.46%
2024/05/20 12:52:41 REST Requests per second: 1258.262891
--- PASS: TestComparePerformance (80.51s)
    --- PASS: TestComparePerformance/gRPC_Performance (1.03s)
    --- PASS: TestComparePerformance/REST_Performance (79.47s)
PASS
ok  	rest_vs_grpc/cmd/bench	80.810s

```