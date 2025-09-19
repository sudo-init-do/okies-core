run:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f api
	
migrate:
	goose -dir ./migrations postgres "host=$$DB_HOST user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "host=$$DB_HOST user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" down
