#!/usr/bin/env bash

function installfluent {
    if systemctl status fluent-bit >/dev/null 2>&1; then
        echo "Fluent Bit is already installed and running. Stopping..."
        sudo systemctl stop fluent-bit
    fi
    if ! /opt/fluent-bit/bin/fluent-bit --version; then
        echo "Fluent Bit is not installed. Installing..."

        if [ ! -f /usr/share/keyrings/fluentbit-keyring.gpg ]; then
            mkdir -p /usr/share/keyrings
        fi
        if [ ! -f /usr/share/keyrings/fluentbit-keyring.gpg ]; then
            curl https://packages.fluentbit.io/fluentbit.key | gpg --dearmor >/usr/share/keyrings/fluentbit-keyring.gpg
            CODENAME=$(lsb_release -cs)
            echo "deb [signed-by=/usr/share/keyrings/fluentbit-keyring.gpg] https://packages.fluentbit.io/ubuntu/${CODENAME} ${CODENAME} main" >>/etc/apt/sources.list
        fi

        sudo apt-get update
        sudo apt-get install -y fluent-bit
    fi

    sudo rm -f /etc/fluent-bit/fluent-bit.conf
    sudo cp /vagrant/obs/fluent-bit.conf /etc/fluent-bit/fluent-bit.conf

    sudo systemctl enable fluent-bit
    sudo systemctl start fluent-bit
    echo "Fluent Bit installed and running."
}
