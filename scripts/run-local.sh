#!/usr/bin/env bash
set -euo pipefail
CONTENT_MODERATION_HTTP_ADDR="${CONTENT_MODERATION_HTTP_ADDR:-:8083}" go run ./cmd/content-moderation-service
