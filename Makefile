.PHONY: run build test clean

run:
	go run cmd/api/main.go

build:
	go build -o bin/billing-engine cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

migrate-up:
	# Add your migration command here

migrate-down:
	# Add your migration command here