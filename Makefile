postgres:
	docker run -d \
		--name postgres \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=password \
			postgres:15.2-alpine

postgres_create:
	PGPASSWORD=password psql -h localhost -p 5432 -d postgres -U postgres -f examples/basic/create.sql

postgres_shell:
	PGPASSWORD=password psql -h localhost -p 5432 -d postgres -U postgres

cockroach:
	docker run -d \
		--name cockroach \
		-p 26257:26257 \
		-p 8080:8080 \
			cockroachdb/cockroach:v23.1.5 \
				start-single-node \
				--insecure

cockroach_create:
	cockroach sql --url "postgres://root@localhost:26257/defaultdb?sslmode=disable" < examples/basic/create.sql

cockroach_shell:
	cockroach sql --url "postgres://root@localhost:26257/defaultdb?sslmode=disable"

shift:
	go run shift.go -c examples/basic/config.yaml

clean:
	docker ps -aq | xargs docker rm -f