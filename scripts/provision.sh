#!/usr/bin/env bash

echo "deb http://ftp.debian.org/debian jessie main contrib" >> /etc/apt/sources.list
apt-get update
apt-get -y upgrade
apt-get -y install curl git build-essential vim jq virtualbox-guest-dkms

# Install and configure go
cd /opt
curl -s -O https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz
tar xzf go1.9.1.linux-amd64.tar.gz

# Install and configure ipfs
curl -s -O https://dist.ipfs.io/go-ipfs/v0.4.11/go-ipfs_v0.4.11_linux-amd64.tar.gz
tar xzf go-ipfs_v0.4.11_linux-amd64.tar.gz
mv go-ipfs/ipfs /usr/local/bin/ipfs

cat <<EOF > /etc/systemd/system/ipfs.service
[Unit]
Description=IPFS daemon

[Service]
Environment="IPFS_PATH=/data/ipfs"
ExecStart=/usr/local/bin/ipfs daemon
Restart=on-failure

[Install]
WantedBy=default.target
EOF
systemctl daemon-reload
systemctl enable ipfs
systemctl restart ipfs



