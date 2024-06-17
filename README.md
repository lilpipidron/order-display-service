# order-display-service

## Test Results By Vegeta

```yaml
Requests:
  - Total: 6000
  - Rate: 50.01
  - Throughput: 50.01

Duration:
  - Total: 2m0s
  - Attack: 2m0s
  - Wait: 3.586ms

Latencies:
  - Min: 1.698ms
  - Mean: 3.035ms
  - 50: 3.089ms
  - 90: 3.52ms
  - 95: 3.695ms
  - 99: 4.067ms
  - Max: 6.181ms

Bytes In:
  - Total: 14310000
  - Mean: 2385.00

Bytes Out:
  - Total: 0
  - Mean: 0.00

Success:
  - Ratio: 100.00%

Status Codes:
  - 200: 6000

Error Set:
```
## Test Results By WRK

```yaml
Threads: 12
Connections: 400

Thread Stats:
  - Avg Latency: 30.73ms
  - Stdev Latency: 7.23ms
  - Max Latency: 247.04ms
  - +/- Stdev Latency: 79.68%
  - Avg Req/Sec: 1.08k
  - Stdev Req/Sec: 169.64
  - Max Req/Sec: 1.97k
  - +/- Stdev Req/Sec: 68.25%

Total Requests: 387590
Total Duration: 30.06s
Total Data Read: 0.91GB

Requests/Sec: 12893.63
Transfer/Sec: 31.09MB
```
