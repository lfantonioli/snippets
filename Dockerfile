FROM ubuntu:22.04 AS builder

# Set the working directory
WORKDIR /app

# Install the necessary packages to build the Go binary
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    golang-go \
    gcc \
    && rm -rf /var/lib/apt/lists/* 


# Copy the source code into the container
COPY . .

# Build the Go binary with flags to enable Go modules and disable CGO
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o app .

# Start from a clean Ubuntu image
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates 

# Set the working directory and the user
WORKDIR /app
# create user and group called appuser and set its UID and GID
RUN groupadd --gid 1000 appuser \
    && useradd --uid 1000 --gid appuser --shell /bin/bash appuser
USER appuser

# Copy the binary from the builder image
COPY --from=builder /app/app .
COPY --from=builder /app/openapi.yaml .

# Expose the port used by the application
EXPOSE 9095

# Start the application
CMD ["./app"]
