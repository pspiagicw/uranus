compile: cmd/uranus/main.go
	go build -o uranus cmd/uranus/main.go
format:
	go fmt ./...
