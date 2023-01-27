#!/usr/bin/env bash

cd $(dirname $0)

docker compose stop

# # postgres コンテナとイメージの削除
# docker container rm -f db_pg
# docker image rm jump_in_pg -f

docker compose \
-f ./docker-compose.yaml \
up --build -d --remove-orphans

echo 'DBをDDLで初期化する場合は、docker imageを削除する必要があります'
echo '詳細はこのシェルの中身参照'
