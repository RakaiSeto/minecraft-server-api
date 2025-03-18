############################
# STEP 1 build executable binary
############################
FROM golang:1.23-alpine3.20 AS builder
RUN apk update && apk add --no-cache gcc musl-dev gcompat
WORKDIR /mengkrep
COPY . .

ENV GOCACHE=/root/.cache/go-build
# Fetch dependencies.
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download
# Build the binary.
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /app/mengkrep

#############################
## STEP 2 build a smaller image
#############################
FROM alpine:3.20
RUN apk update && apk add --no-cache ffmpeg
WORKDIR /app

RUN mkdir -p /root/.ssh && touch /root/.ssh/known_hosts

# Copy compiled from builder.
COPY --from=builder /app/mengkrep /app/mengkrep

COPY static /app/static

COPY .env .env

# Run the binary.
ENTRYPOINT ["/app/mengkrep"]