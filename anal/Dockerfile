FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build all binaries
RUN go build -o /bin/calculate_sums ./cmd/calculate_sums/
RUN go build -o /bin/generate_pyramid ./cmd/generate_pyramid/
RUN go build -o /bin/analyze_titles ./cmd/analyze_titles/
RUN go build -o /bin/run_all ./cmd/run_all/
RUN go build -o /bin/upload_analysis ./cmd/upload_analysis/

# Final stage
FROM alpine:latest

# Install ca-certificates for any external API calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /bin/calculate_sums /bin/
COPY --from=builder /bin/generate_pyramid /bin/
COPY --from=builder /bin/analyze_titles /bin/
COPY --from=builder /bin/run_all /bin/
COPY --from=builder /bin/upload_analysis /bin/

# Create output directory
RUN mkdir -p /app/output

# Default command
CMD ["/bin/run_all", "-data", "/data", "-output", "/app/output", "-workers", "10"]