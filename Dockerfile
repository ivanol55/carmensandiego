FROM golang:1.20.4-alpine3.18 AS builder

# Copy the repository into the working directory
WORKDIR /carmensandiego
COPY . .

# Install dependencies
RUN go install

# Compile the helper executable
RUN go build helper.go

# Remove .git files to avoid repository disclosure
RUN rm -rf src/ .git .gitignore README.md docker-compose.yaml Dockerfile  go.mod go.sum LICENSE

# Switch to Alpine image as a final base
FROM alpine:3.18.0 AS helper

# Copy built artifacts from builder container
WORKDIR /carmensandiego
COPY --from=builder /carmensandiego ./

ENTRYPOINT ["./helper"]