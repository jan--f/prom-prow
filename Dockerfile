FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot ./cmd/bot

# Final stage
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /build/bot /bot

USER nonroot:nonroot

ENTRYPOINT ["/bot"]
