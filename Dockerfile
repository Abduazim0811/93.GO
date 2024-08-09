FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .

CMD ["./main"]
