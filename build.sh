#!/usr/bin/env bash
GIT_COMMIT=$(git describe --tags)
DEFAULT_SERVER="https://lemon.everyclass.xyz"
ENABLE_METRICS="true"
LANGUAGE="zh"
go build -ldflags "-X main.ravenDSN=$SENTRY_DSN -X main.gitRevision=$GIT_COMMIT -X main.enableMetrics=$ENABLE_METRICS -X main.enableGlobalReport=true -X main.defaultServer=$DEFAULT_SERVER -X main.lang=$LANGUAGE" -o ./bin/lemon ./lemon