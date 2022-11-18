#!/bin/bash
if [ ! $@ ]; then
name='student-container'
else
name=$@
fi
docker rmi `docker ps -aq -f name=$name`
set -a
source .env

docker run --rm \
--network student-networks \
-e LOG_LEVEL=$LOG_LEVEL \
-e SERVICE_NAME=$SERVICE_NAME \
-e ENVIRONMENT=$ENVIRONMENT \
-e TRACER_ID=$TRACER_ID \
-e TRACER_EXPORTER_URL=$TRACER_EXPORTER_URL \
-e DATABSE_NAME=$DATABSE_NAME \
-e DATABASE_URL=$DATABASE_URL \
-e HTTP_PORT=$HTTP_PORT \
-e DATABASE_MIGRATION_PATH=$DATABASE_MIGRATION_PATH \
-e HTTP_DIAL_TIMEOUT=$HTTP_DIAL_TIMEOUT \
-p 15343:15343 \
--name $name student-img