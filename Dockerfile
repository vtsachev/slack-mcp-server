FROM golang:1.23 AS build

ENV CGO_ENABLED=0
ENV GOTOOLCHAIN=local
ENV GOCACHE=/go/pkg/mod

RUN apt-get update  \
  && apt-get install -y --no-install-recommends net-tools curl

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . /app

RUN --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -w" -o /go/bin/mcp-server ./cmd/slack-mcp-server

FROM build AS dev

RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/go-delve/delve/cmd/dlv@v1.23.1 && cp /go/bin/dlv /dlv

WORKDIR /app/mcp-server

EXPOSE 3001

CMD ["mcp-server", "--transport", "sse"]

FROM alpine:3.18 AS production

RUN apk add --no-cache ca-certificates

# Create a new group and user
RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

COPY --from=build /go/bin/mcp-server /usr/local/bin/mcp-server

# Ensure the binary is executable by the nonroot user
RUN chmod +x /usr/local/bin/mcp-server

WORKDIR /app

# Change ownership of the /app directory
RUN chown -R nonroot:nonroot /app

# Switch to the nonroot user
USER nonroot

EXPOSE 3001

CMD ["mcp-server", "--transport", "sse"]
