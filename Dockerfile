# Production image for running the activecookie CLI.
#
# Build:
#   docker build -t activecookie .
#
# Push to Docker Hub:
#   docker tag activecookie <username>/activecookie:latest
#   docker push <username>/activecookie:latest
#
# Run:
#   docker run --rm -v "$(pwd):/data" activecookie \
#     -f /data/cookie_log.csv -d 2018-12-09

FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /activecookie ./cmd/activecookie

FROM scratch

COPY --from=builder /activecookie /activecookie

ENTRYPOINT ["/activecookie"]
