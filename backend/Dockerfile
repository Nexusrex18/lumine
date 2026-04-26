# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy module files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source
COPY . .

# Build only the backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./backend

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.20

# Non-root user for security
RUN addgroup -S app && adduser -S app -G app
USER app

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

ENV PORT=8080
ENV CORS_ALLOWED_ORIGINS=*

ENTRYPOINT ["./server"]
