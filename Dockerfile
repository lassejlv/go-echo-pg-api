FROM golang:1.23.6-alpine AS builder
 
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
RUN go run github.com/steebchen/prisma-client-go prefetch
COPY ./ ./
RUN go run github.com/steebchen/prisma-client-go generate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .
FROM scratch
COPY --from=builder /workspace/app /app
ENTRYPOINT ["/app"]
 