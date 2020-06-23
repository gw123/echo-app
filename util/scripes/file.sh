#!/usr/bin/env bash
APP=file
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

<<<<<<< HEAD
<<<<<<< HEAD:util/scripes/file.sh
./echoapp file --config=config.yaml &
=======
./echoapp comment --config=config.yaml &
>>>>>>> develop:util/scripes/run.sh
=======
./echoapp file --config=config.yaml &
>>>>>>> develop

child=$!
wait "$child"

