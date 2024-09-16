FROM golang:1.22.3 AS builder

WORKDIR /usr/local/src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app cmd/main.go

FROM ubuntu:latest AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/migrations /migrations

EXPOSE 8080
CMD ["/app"]
