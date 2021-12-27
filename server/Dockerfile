FROM golang:latest AS builder
LABEL maintainer="al.suhodubenko@yandex.ru"

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tmp-server ./cmd/main.go

FROM alpine:latest
COPY --from=builder /app/tmp-server .
CMD ["./tmp-server"]