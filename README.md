# wbtech

1) Docker
```sh
docker compose build
docker compose up -d
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

## TODO:
- instruction to start 
- Dockerfile (compose) - DONE
- slog 
- test && mocks
- config (yaml) - DONE
- clear cache (control memory) - DONE
- html template
- github actions - DONE