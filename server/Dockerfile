FROM golang:1.21 as builder

WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/server

FROM debian:buster-slim
COPY --from=builder /app/.env app/.env
COPY --from=builder /app/bin/server app/server

EXPOSE 3000
CMD cd /app && ./server
