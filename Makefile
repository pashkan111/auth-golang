
start: 
	docker-compose up -d

stop: 
	docker-compose down

test:
	go test ./src/tests
