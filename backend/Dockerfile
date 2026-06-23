FROM golang:1.25.1-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /out/api ./cmd/api

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates && adduser -D -g '' appuser

COPY --from=build /out/api /app/api

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/api"]
