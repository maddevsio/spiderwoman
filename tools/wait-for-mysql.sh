#!/usr/bin/env bash

# wrapper script to check if mysql is up and then start service
# after 3 successful attempets or exit after timeout.

HOST=$1
shift
TIMEOUT=$1
shift
CMD=$@

start_sec=$(date +%s)
end_sec=$(date +%s)
count=0

while [ $((end_sec - start_sec)) -lt $TIMEOUT ]; do
    (echo > /dev/tcp/$HOST/3306) >/dev/null 2>&1
    result=$?
    if [[ $result -eq 0 ]]; then
        let count+=1
        if [[ $count -ge 3 ]]; then
            >&2 echo "Mysql($HOST) is up - starting service"
            exec $CMD
            break
        fi
    else
        let count=0
    fi
    >&2 echo "Waiting if Mysql($HOST) is up.."
    sleep 5
    end_sec=$(date +%s)
done

if [[ $result -ne 0 ]]; then
    >&2 echo "Mysql($HOST) did not start after $TIMEOUT sec"
    exit $result
fi
