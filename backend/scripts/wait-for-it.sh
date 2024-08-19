#!/bin/sh

TIMEOUT=15
HOST=$1
PORT=$2
shift 2

for i in $(seq 1 $TIMEOUT); do
    nc -z "$HOST" "$PORT" > /dev/null 2>&1
    result=$?
    if [ $result -eq 0 ]; then
        echo "$HOST:$PORT is available after $i seconds"
        exec "$@"
        exit 0
    fi
    sleep 1
done

echo "timeout occurred after waiting $TIMEOUT seconds for $HOST:$PORT"
exit 1
