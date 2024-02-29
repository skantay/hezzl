FROM golang:1.20 AS build_base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/service-1 ./cmd/api/main.go

FROM alpine:3.9

COPY --from=build_base /app/service-1 /app/service-1

COPY config /app/config
COPY migrations /app/migrations

EXPOSE 8080

WORKDIR /app

CMD ["/app/service-1"]
