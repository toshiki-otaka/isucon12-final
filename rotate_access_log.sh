#!/bin/sh

set -eu

# log rotation
sudo mv /var/log/nginx/access.log /var/log/nginx/access.log.`date +%Y%m%d%H%M%S`
sudo nginx -s reopen

echo "success!"
