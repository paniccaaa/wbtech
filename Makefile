PG="postgres://wbuser:wbpassword@localhost:5435/postgres?sslmode=disable"

goose-up:
	@goose -dir ./migrations postgres $(PG) up

goose-down:
	@goose -dir ./migrations postgres $(PG) down

connect:
	@psql $(PG)
