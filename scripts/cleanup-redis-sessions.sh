#!/bin/bash

# Script to clean up Redis session keys
echo "Cleaning up Redis session keys..."

# Check if redis-cli is available
if ! command -v redis-cli &> /dev/null; then
    echo "redis-cli not found. Please install redis-tools or ensure Redis CLI is available."
    echo "You can manually connect to Redis and run: EVAL \"return redis.call('del', unpack(redis.call('keys', 'session:*')))\" 0"
    exit 1
fi

# Count existing session keys
SESSION_COUNT=$(redis-cli --scan --pattern "session:*" | wc -l)
echo "Found $SESSION_COUNT session keys to delete"

if [ "$SESSION_COUNT" -gt 0 ]; then
    # Delete all session keys
    redis-cli --scan --pattern "session:*" | xargs redis-cli del
    echo "Deleted $SESSION_COUNT session keys"
else
    echo "No session keys found to delete"
fi

echo "Redis cleanup completed!"
