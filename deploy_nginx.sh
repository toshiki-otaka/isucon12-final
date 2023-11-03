#!/bin/sh

set -eu

# nginx
sudo cp -r ./nginx/nginx.conf /etc/nginx
sudo cp -r ./nginx/isuconquest.conf /etc/nginx/sites-available
sudo nginx -s reload

echo "success!"
