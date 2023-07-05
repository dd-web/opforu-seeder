build:
	@go build -o bin/seeder.exe

run: build
	@./bin/seeder.exe

test:
	@go test -v ./...