#!/bin/bash

kill $(pgrep -f nginx)

PORT=8001 ./volume /tmp/volume1/ &
PORT=8002 ./volume /tmp/volume2/ &
PORT=8003 ./volume /tmp/volume3/ &
PORT=8004 ./volume /tmp/volume4/ &
PORT=8005 ./volume /tmp/volume5/ &
PORT=8006 ./volume /tmp/volume6/ &
PORT=8007 ./volume /tmp/volume7/ &
PORT=8008 ./volume /tmp/volume8/ &


./steres -port 8000 -db \
volumes  localhost:8001, \
localhost:8002, \
localhost:8003, \
localhost:8004, \
localhost:8005, \
localhost:8006, \
localhost:8007, \
localhost:8008, \
-db /db serve
