# ==============================================================================
# Build Stage
# ==============================================================================
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git curl nodejs npm

# Install Go tools
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/templui/templui/cmd/templui@latest

# Install Tailwind CSS v4 CLI globally via npm for architecture compatibility (arm64/amd64)
RUN npm install -g @tailwindcss/cli tailwindcss
ENV NODE_PATH=/usr/local/lib/node_modules

WORKDIR /app

# Copy go module files first for better Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Download and add templui components on build
RUN templui add "*"

# Compile templ templates and tailwind styles
RUN templ generate
RUN tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css

# Build the optimized production binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/main.go

# ==============================================================================
# Run Stage
# ==============================================================================
FROM alpine:latest

# Install certificates for SSL calls
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the compiled binary and configuration
COPY --from=builder /out/server /app/server
COPY --from=builder /app/config.yaml /app/config.yaml
# Copy static assets (including generated CSS and templui JS files)
COPY --from=builder /app/assets /app/assets

EXPOSE 8080

ENTRYPOINT ["/app/server"]
