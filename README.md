# wbtech L0

1) Docker
```sh
docker compose build
docker compose up -d
```

- other way:
```sh
CONFIG_PATH="./config/dev.yaml" go run cmd/wbtech/main.go -env=local
```
2) Migrations
- Install -> https://github.com/pressly/goose
```sh
make goose-up
```
- Result:
```sh
2024/11/21 12:55:30 OK   20241111162423_orders_table.sql (6.95ms)
2024/11/21 12:55:30 goose: successfully migrated database to version: 20241111162423
```

- To check DB:
```sh
make connect
```
3) Endpoints

- GetOrderByID:

For example: 
```sh
curl http://localhost:8089/order/b563feb7b2b84b6test0 
```

- KafkaUI: [localhost:8080](http://localhost:8080/)

4) Test service 
- wrk test:
```sh
wrk -t12 -c400 -d30s http://localhost:8089/order/b563feb7b2b84b6test0 
``` 

- Result:
```sh
Running 30s test @ http://localhost:8089/order/b563feb7b2b84b6test0
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    12.62ms    8.00ms 135.94ms   75.68%
    Req/Sec     2.76k   416.57     3.70k    75.83%
  988961 requests in 30.07s, 623.42MB read
Requests/sec:  32887.88
Transfer/sec:     20.73MB
```