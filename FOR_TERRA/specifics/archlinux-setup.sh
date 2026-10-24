#! /bin/sh

mkdir -p ~/.config/pkgeep
cp config.yaml ~/.config/pkgeep/

pacman -Syu
~/pkgeep setup

echo 'PATH="$PATH:/root/"' >> ~/.bashrc

echo 'All setup, run `pkgeep help` to get started!'

exec bash
