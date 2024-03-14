FROM golang:latest as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o uranus cmd/uranus/main.go

FROM scratch

ENV USER=uranus

COPY --from=builder /usr/src/app/uranus /

CMD ["/uranus"]
