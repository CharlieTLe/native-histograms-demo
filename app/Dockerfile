FROM golang:1.23-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o /demo-app .

FROM alpine:3.20
COPY --from=builder /demo-app /demo-app
ENTRYPOINT ["/demo-app"]
