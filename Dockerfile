FROM golang:1.23.6-alpine AS builder
 
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .
FROM scratch
COPY --from=builder /workspace/app /app
ENTRYPOINT ["/app"]
 