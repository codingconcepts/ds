validate_version:
ifndef VERSION
	$(error VERSION is undefined)
endif

postgres:
	docker run -d \
		--name postgres \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=password \
			postgres:15.2-alpine

postgres_create:
	PGPASSWORD=password psql -h localhost -p 5432 -d postgres -U postgres -f examples/basic/create.sql

postgres_insert:
	PGPASSWORD=password psql -h localhost -p 5432 -d postgres -U postgres -f examples/basic/insert.sql

postgres_update:
	PGPASSWORD=password psql -h localhost -U postgres -c 'UPDATE person SET full_name = upper(full_name)'

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

dshift:
	go run ds.go insert --config examples/basic/config.yaml

verify:
	molt verify \
		--source 'postgres://postgres:password@localhost:5432/postgres?sslmode=disable' \
		--target 'postgresql://root@localhost:26257/defaultdb?sslmode=disable'

test:
	go test ./... -v -cover

cover:
	go test -v -coverpkg=./... -coverprofile=coverage.out ./... -count=1
	go tool cover -html coverage.out

release: validate_version
	GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o ds ;\
	tar -zcvf ./releases/ds_${VERSION}_linux.tar.gz ./ds ;\

	GOOS=darwin go build -ldflags "-X main.version=${VERSION}" -o ds ;\
	tar -zcvf ./releases/ds_${VERSION}_macOS.tar.gz ./ds ;\

	GOOS=windows go build -ldflags "-X main.version=${VERSION}" -o ds ;\
	tar -zcvf ./releases/ds_${VERSION}_windows.tar.gz ./ds ;\

	rm ./ds

clean:
	docker ps -aq | xargs docker rm -f