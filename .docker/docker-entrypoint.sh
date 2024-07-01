#!/usr/bin/env sh
set -ue

set -x ## debug trace

case $1 in
  test-integ)
    gocov test -v -timeout 200s ./... | gocov report
    exit
    ;;
  test-gocov)
    # -race not working on alpine: https://github.com/golang/go/issues/14481
    # CGO_ENABLED=1 go build -race -o study-manager-test .
    #gocov test . | gocov-xml > coverage.xml
    gocov test -tags=unit ./... | gocov report
    exit
    ;;
  test-swagger)
    exec test -s /github.com/stokkelol/intospaceapi/static/api.json
    ;;
  linter)
    exec golangci-lint run -v ./...
    ;;
  air)
    shift
    cd /app/cmd
    exec air "$@"
    ;;
  delve)
    shift
    cd /app
    CGO_ENABLED=0 go build -x -gcflags="-N -l" -o intospaceapi-delve .
    chmod +x ./intospaceapi-delve
    dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient --log exec ./intospaceapi-delve -- "$@"
    exit
    ;;
  test-cli)
    exec go test ./test-cli
    exit
    ;;
esac

exec "$@"