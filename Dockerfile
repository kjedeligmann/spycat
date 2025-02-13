FROM golang:1.23-bookworm AS base

# Development stage
# =============================================================================
# Create a development stage based on the "base" image
FROM base AS development

WORKDIR /app

# Install the air CLI for auto-reloading
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

# Start air for live reloading
CMD ["air"]

# Builder stage
# =============================================================================
# Create a builder stage based on the "base" image
FROM base AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
# Turn off CGO to ensure static binaries
RUN CGO_ENABLED=0 go build -o spycat

# Production stage
# =============================================================================
# Create a production stage to run the application binary
FROM scratch AS production

WORKDIR /prod

# Copy binary from builder stage
COPY --from=builder /build/spycat ./

# Document the port that may need to be published
EXPOSE 8080

# Start the application
CMD ["/prod/spycat"]

