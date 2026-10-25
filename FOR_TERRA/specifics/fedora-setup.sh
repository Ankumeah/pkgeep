#! /bin/sh

mkdir -p ~/.config/pkgeep
cp config.yaml ~/.config/pkgeep/

dnf install libdnf5-plugin-actions
~/pkgeep setup

echo 'PATH="/root/:$PATH"' >> ~/.bashrc

echo 'All setup, run `pkgeep help` to get started!'

exec bash
