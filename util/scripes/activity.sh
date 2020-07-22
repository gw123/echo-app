#!/usr/bin/env bash
APP=site
term() {
    echo "Caught SIGTERM signal!"
    kill -TERM "$child" 2>/dev/null
}

if [ "$1" == "update-config" ]; then
  etcdctl ${AUTH} ${ENDPOINTS} get /prod/${APP}/config.yaml > ./config.new.yaml
fi

trap _term SIGTERM

if [ -f upload ]; then
    mv upload echoapp
    chmod +x echoapp
fi

./echoapp activity --config=config.yaml &

child=$!
wait "$child"

