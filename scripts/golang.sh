#!/bin/sh

VERSION="$1"

install_go()
{
  mkdir -p ~/go/bin
  test -f /usr/bin/sudo && sudo rm -rf /usr/local/go ~/.cache/go-build/* ~/go/* /tmp/go${VERSION}.linux-amd64.tar.gz || rm -rf /usr/local/go ~/.cache/go-build/* ~/go/* /tmp/go${VERSION}.linux-amd64.tar.gz
  wget https://go.dev/dl/go${VERSION}.linux-amd64.tar.gz -P /tmp

  test -f /usr/bin/sudo && sudo tar -C /usr/local -xzf /tmp/go${VERSION}.linux-amd64.tar.gz || tar -C /usr/local -xzf /tmp/go${VERSION}.linux-amd64.tar.gz
  rm -rf /tmp/go${VERSION}.linux-amd64.tar.gz
  grep /usr/local/go/bin ~/.bashrc
  if [ $? -eq 1 ]
  then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
    echo 'export GO111MODULE=on' >> ~/.bashrc
    . ~/.bashrc
  fi
  which go
}

if [ -n "$VERSION" ]
then
  which go || install_go
  go version |grep "$VERSION" || install_go
else
  echo "golang version is missing"
  exit 1
fi
