#!/bin/sh
set -e

echo "Creating volumes"

sudo mkdir -p /srv/db /srv/grafana /srv/victoria-metrics /srv/victoria-logs /srv/vmagent

sudo chown -R 999:999 /srv/db
sudo chown -R 472:472 /srv/grafana

sudo chmod -R 700 /srv/db
sudo chmod -R 755 /srv/grafana

echo "Volumes created"
