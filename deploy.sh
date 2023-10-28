#!/bin/sh

set -eu

APP_HOME=/home/isucon/webapp

# log rotation
sudo mv /var/log/nginx/access.log /var/log/nginx/access.log.`date +%Y%m%d%H%M%S`
sudo nginx -s reopen

if [ -e /tmp/mysql-slow.log ]; then
  sudo cp /tmp/mysql-slow.log /var/log/mysql/slowquery.log.`date +%Y%m%d%H%M%S`
fi

# NOTE: chown root /tmp/mysql-slow.logしたら動いた
sudo truncate -s 0 /tmp/mysql-slow.log

BRANCH=${1:-main}
git fetch origin $BRANCH
git switch $BRANCH
git pull origin $BRANCH

# app
cd webapp/go
/home/isucon/local/golang/bin/go build -o isuconquest .
mv isuconquest ${APP_HOME}/go/isuconquest
cd -

systemctl restart isuconquest.go.service

echo "success!"
