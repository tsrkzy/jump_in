#!/bin/bash
echo 'docker-entrypoint.sh'

cd /var/data/ddl

for f in *.sql
do
 echo ${f}
done

echo '>> docker-entrypoint.sh done.'
