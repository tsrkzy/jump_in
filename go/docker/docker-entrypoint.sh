#!/bin/bash
echo 'docker-entrypoint.sh'


# nginx conf
set +e
rm -rf /etc/nginx/nginx.conf
rm -rf /etc/nginx/conf.d/default.conf
ln -s /var/www/docker/nginx.conf /etc/nginx/nginx.conf
set -e

# go environment file
# @SEE jump_in/go/docker/go_echo.service
envsubst < /etc/default/_go_echo > /etc/default/go_echo
rm -rf /etc/default/_go_echo

# service
systemctl enable nginx.service
systemctl start nginx.service

systemctl enable go_echo.service
systemctl start go_echo.service

systemctl restart rsyslog

# agettyの停止 (docker-composeでprivilegeを指定 && Mac だと発生することを確認。CPUを食いつぶす)
systemctl disable getty@tty1.service
systemctl stop getty@tty1.service
systemctl disable system-getty.slice
systemctl stop system-getty.slice
systemctl disable getty.target
systemctl stop getty.target

echo '>> docker-entrypoint.sh done.'
