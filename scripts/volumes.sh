#!/bin/sh
set -e

sudo mkdir -p /srv/grafana /srv/loki /srv/db

sudo chown -R 472:472 /srv/grafana
sudo chown -R 10001:10001 /srv/loki
sudo chown -R 999:999 /srv/db

sudo chmod -R 755 /srv/grafana
sudo chmod -R 755 /srv/loki
sudo chmod -R 700 /srv/db

echo "Volumes created"
