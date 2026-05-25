FROM golang:1.26-alpine AS build

WORKDIR /src

COPY Source/go.mod Source/go.sum ./
RUN go mod download

COPY Source/ ./

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/server ./cmd/server

FROM alpine:3.19

RUN adduser -D -u 1000 app

WORKDIR /app

COPY --from=build /out/server /app/server
COPY Source/mocks.json /app/mocks.json

RUN chown -R app:app /app

USER app

ENV MOCKS_FILE=/app/mocks.json
ENV PORT=8080

EXPOSE 8080

CMD ["/app/server"]