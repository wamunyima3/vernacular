# The API and its command-line tools in one image (P1-OPS-01). Build from the repo root:
#   docker build -f apps/api/Dockerfile -t vernacular-api .
FROM golang:1.26-alpine AS build
WORKDIR /src/apps/api
COPY apps/api/go.mod apps/api/go.sum ./
COPY apps/api/tools/go.mod apps/api/tools/go.sum tools/
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod download -modfile=tools/go.mod
COPY apps/api ./
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/ ./cmd/api ./cmd/worker ./cmd/ingest ./cmd/admin && \
    go build -trimpath -ldflags="-s -w" -modfile=tools/go.mod -o /out/ \
      github.com/pressly/goose/v3/cmd/goose github.com/riverqueue/river/cmd/river

FROM alpine:3.22
# ca-certificates: the worker and ingest call HTTPS services.
# /data/blobs is owned by the app user, so a fresh volume mounted there is too.
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H -u 10001 vernacular && \
    mkdir -p /data/blobs && chown vernacular /data/blobs
COPY --from=build /out/ /usr/local/bin/
COPY apps/api/db/migrations /app/db/migrations
COPY apps/api/deploy/migrate.sh /usr/local/bin/migrate
COPY docs/style-guides /app/docs/style-guides
COPY data /app/data
WORKDIR /app
ENV HTTP_ADDR=:8080 STYLE_GUIDES_DIR=/app/docs/style-guides
USER vernacular
EXPOSE 8080
CMD ["api"]
