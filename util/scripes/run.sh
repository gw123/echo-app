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

<<<<<<< HEAD
./echoapp comment --config=config.yaml &
=======
<<<<<<< HEAD:util/scripes/file.sh
./echoapp file --config=config.yaml &
=======
./echoapp comment --config=config.yaml &
>>>>>>> develop:util/scripes/run.sh
>>>>>>> develop

child=$!
wait "$child"

