#!/usr/bin/env bash
GIT_COMMIT=$(git describe --tags)
DEFAULT_SERVER="https://lemon.everyclass.xyz"
ENABLE_METRICS="true"
go build -ldflags "-X main.ravenDSN=$SENTRY_DSN -X main.gitRevision=$GIT_COMMIT -X main.enableMetrics=$ENABLE_METRICS -X main.enableGlobalReport=true -X main.defaultServer=$DEFAULT_SERVER" -o ./bin/lemon ./lemon