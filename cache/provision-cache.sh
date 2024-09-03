#!/usr/bin/env bash

export DEBIAN_FRONTEND=noninteractive

sudo apt-get update
sudo apt-get install -y build-essential memcached

echo "reloading memcached system unit..."
sudo systemctl stop memcached
sudo cp /vagrant/cache/memcached.conf /etc/memcached.conf
sudo systemctl start memcached
echo "reloaded memcached"

if [ ! -f /vagrant/obs/provision-fluentbit.sh ]; then
    echo "provision-fluentbit.sh not found"
    echo "no fluent bit logging will run"
    exit
fi

echo "provision-fluentbit.sh found"
source /vagrant/obs/provision-fluentbit.sh

installfluent 10.0.0.26
