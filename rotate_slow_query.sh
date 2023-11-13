#!/bin/sh

set -eu

if [ -e /tmp/mysql-slow.log ]; then
  sudo cp /tmp/mysql-slow.log /var/log/mysql/slowquery.log.`date +%Y%m%d%H%M%S`
fi

sudo truncate -s 0 /tmp/mysql-slow.log

echo "success!"
