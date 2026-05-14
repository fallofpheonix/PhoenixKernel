run:
	docker compose up -d

stop:
	docker compose down

test:
	pytest tests/

lint:
	ruff check .

build:
	docker compose build

clean:
	docker compose down -v
	find . -type d -name __pycache__ -exec rm -rf {} +

logs:
	docker compose logs -f

