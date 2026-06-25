FROM golang:1.25.7-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY internal/postgres/migrations ./internal/postgres/migrations

RUN go build -o /out/api ./cmd/api
RUN go build -o /out/migrate ./cmd/migrate

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates && adduser -D -g '' appuser

COPY --from=build /out/api /app/api
COPY --from=build /out/migrate /app/migrate
COPY --from=build /src/internal/postgres/migrations /app/migrations

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/api"]
