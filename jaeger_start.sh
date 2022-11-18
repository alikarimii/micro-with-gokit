#!/bin/bash
if [ ! $@ ]; then
name='jaeger-container'
else
name=$@
fi
docker rmi `docker ps -aq -f name=$name`
set -a
source .env

docker run -d \
--network student-networks \
-p 6831:6831/udp \
-p 16686:16686 \
-p 14268:14268 \
--name $name jaegertracing/all-in-one:latest
