#!/usr/bin/env bash
# docker(go_echo)の中にairが入っているので、基本的にリロードは意識しなくて良い
# @REF https://github.com/cosmtrek/air

echo 'go build... '
docker exec -it go_echo go mod tidy
docker exec -it go_echo go build -v -o /usr/local/bin/app
echo 'done!'
