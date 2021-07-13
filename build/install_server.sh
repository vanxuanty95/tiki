#!/bin/bash

PWD=$(dirname "$0")
docker-compose -f $PWD/docker-compose.yml up -d --build

sleep 10

docker cp $PWD/init_database.sql tiki_mysql:/docker-entrypoint-initdb.d/init_database.sql
docker exec -it tiki_mysql psql -d togo -U user -f docker-entrypoint-initdb.d/init_database.sql

docker-compose -f $PWD/docker-compose.yml stop

echo Done!