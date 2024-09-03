#!/usr/bin/env bash

set -e

if [ ! -f /vagrant/obs/provision-fluentbit.sh ]; then
    echo "provision-fluentbit.sh not found"
    exit 1
fi

echo "provision-fluentbit.sh found"
source /vagrant/obs/provision-fluentbit.sh

installfluent 10.0.0.26
