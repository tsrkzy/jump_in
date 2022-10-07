#!/usr/bin/env bash

echo ' - entry-point: go_echo'
docker exec -it go_echo sh /var/www/docker/docker-entrypoint.sh

#echo ' - entry-point: db_pg'
#docker exec -it db_pg sh /var/data/docker/docker-entrypoint.sh
