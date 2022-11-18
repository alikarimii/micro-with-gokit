#!/bin/bash
if [ ! $@ ]; then
name='postgres14-container'
else
name=$@
fi
docker rmi `docker ps -aq -f name=$name`
set -a
source .env

docker run -d \
--network student-networks \
-e STUDENT_LOCAL_DATABASE=$STUDENT_LOCAL_DATABASE \
-e STUDENT_TEST_DATABASE=$STUDENT_TEST_DATABASE \
-e POSTGRES_USER=$POSTGRES_USER \
-e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
-e STUDENT_USERNAME=$STUDENT_USERNAME \
-e STUDENT_PASSWORD=$STUDENT_PASSWORD \
-v "$PWD/internal/infrastructure/postgres/database/setup:/docker-entrypoint-initdb.d" \
--name $name postgres:14.4
