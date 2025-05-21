build:
	docker build -t solar-garlic-mailing-list .

run:
	docker run -p 8080:8080 --env-file .env solar-garlic-mailing-list

fmt:
	go fmt ./...