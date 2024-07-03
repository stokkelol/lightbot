## BASE STAGE
FROM golang:1.22 AS builder

WORKDIR /app

ARG APP_VER
ARG BUILD_DATE

COPY ./go.mod ./go.sum ./
RUN go mod download
RUN go install github.com/air-verse/air

COPY ./ ./
RUN cd cmd && CGO_ENABLED=0 GOOS=linux GOEXPERIMENT=rangefunc go build -o app \
  -ldflags="-X 'internal/service.Version=${APP_VER}' -X 'internal/service.BuildDate=${BUILD_DATE}'"

## RELEASE STAGE
FROM alpine AS release
ARG ENV
ARG PORT
ENV APP_ENV="${ENV}"
WORKDIR /app

## copy needed
COPY --from=builder /app/cmd/app ./
EXPOSE $PORT
CMD ["./app"]

## TEST STAGE
FROM builder AS test
ENV APP_ENV="test"

## setup needed env vars for tests
EXPOSE 9000
WORKDIR /app
COPY ./docker/docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
COPY --from=builder /app/migrations ./migrations
CMD ["./app"]

## RELEASE STAGE
FROM release