build: cmd/uranus/main.go
	go build -o uranus cmd/uranus/main.go
format:
	go fmt ./...
test:
	go test ./...
run:
	go run cmd/uranus/main.go
