worker-start:
	docker compose --profile worker up worker -d --build

worker-stop:
	docker compose --profile worker down

worker-logs:
	docker compose logs worker -f
