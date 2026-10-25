#! /bin/sh

mkdir -p /root/.config/pkgeep
cp config.yaml /root/.config/pkgeep/

apt update
/root/pkgeep setup

echo 'PATH="/root/:$PATH"' >> /root/.bashrc

echo 'All setup, run `pkgeep help` to get started!'

exec bash
