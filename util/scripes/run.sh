#!/usr/bin/env bash
term() {
    echo "Caught SIGTERM signal!"
    kill -TERM "$child" 2>/dev/null
}

trap _term SIGTERM

if [ -f upload ]; then
    mv upload echoapp
    chmod +x echoapp
fi

./echoapp comment --config=config.yaml &

child=$!
wait "$child"
