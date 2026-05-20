FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/server

FROM alpine:3.19
RUN adduser -D -u 1000 app && mkdir -p /data /seed && chown -R app:app /data /seed
COPY --from=build /out/app /app
COPY mocks.json /seed/mocks.json
ENV MOCKS_FILE=/data/mocks.json \
    SEED_FILE=/seed/mocks.json \
    PORT=8080
USER app
EXPOSE 8080
ENTRYPOINT ["/app"]
