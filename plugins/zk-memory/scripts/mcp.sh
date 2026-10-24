#!/bin/sh
# Runs the zk MCP server against the memory vault.
. "$(dirname "$0")/vault.sh" || exit 1

ZK_VAULT="$V" exec zk mcp
