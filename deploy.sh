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
cp ./webapp/sql/init.sh ${APP_HOME}/sql/init.sh
cp ./webapp/sql/3_schema_exclude_user_presents.sql ${APP_HOME}/sql/3_schema_exclude_user_presents.sql
cp ./webapp/sql/4_alldata_exclude_user_presents* ${APP_HOME}/sql/
cp ./webapp/sql/5_user_presents_not_receive_data* ${APP_HOME}/sql/
cp ./provisioning/packer/ansible/roles/xbuild/files/home/isucon/env /home/isucon/env

sudo systemctl restart isuconquest.go.service

echo "success!"
