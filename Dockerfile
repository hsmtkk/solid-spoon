FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY web.go .
RUN CGO_ENABLED=0 go build -o web

FROM AS runtime
COPY --from=builder /app/web /app/web
ENTRYPOINT ["/app/web"]
