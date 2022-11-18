#!/bin/bash
if [ ! $@ ]; then
name='student-img'
else
name=$@
fi
docker build -t $name \
-f Dockerfile .