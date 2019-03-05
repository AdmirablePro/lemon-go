#!/usr/bin/env bash
GIT_COMMIT=$(git describe --tags)
go build -ldflags "-X main.ravenDSN=$SENTRY_DSN -X main.gitRevision=$GIT_COMMIT" -o ./bin/lemon ./lemon