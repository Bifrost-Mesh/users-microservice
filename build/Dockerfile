# syntax=docker/dockerfile:1

# Stage 1 : Build the binary.
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN --mount=type=bind,src=go.mod,target=go.mod \
  --mount=type=bind,src=go.sum,target=go.sum \
  --mount=type=cache,target=/go/pkg/mod  \
  go mod download

RUN --mount=type=bind,src=.,target=. \
  --mount=type=cache,target=/go/pkg/mod  \
  CGO_ENABLED=0 GOOS=linux go build -o /users-microservice ./cmd

# Stage 2 : Package the binary.
FROM scratch AS packager

LABEL org.opencontainers.image.authors="archismanmridha12345@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/Bifrost-Mesh/users-microservice"
LABEL org.opencontainers.image.license="MIT"

WORKDIR /

COPY --from=builder /users-microservice /users-microservice

USER nonroot:nonroot
EXPOSE 4000

ENTRYPOINT ["/users-microservice"]
