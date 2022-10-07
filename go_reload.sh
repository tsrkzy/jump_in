#!/usr/bin/env bash
# docker用
echo 'go build... '
docker exec -it go_echo go mod tidy
docker exec -it go_echo go build -v -o /usr/local/bin/app
echo 'done!'
printf 'systemctl restart ... '
docker exec -it go_echo systemctl restart go_echo
echo 'done!'
