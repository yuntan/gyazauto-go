#!/bin/bash

build () {
  os=$1
  arch=$2
  echo "for ${os} ${arch}bit"

  bin=bin/gyazauto-$(git describe)-$os
  if [[ $os == "windows" ]]; then
    bin=$bin.exe
  fi
  arch_=amd64
  if [[ $arch == "32" ]]; then
    arch_="386"
  fi
  GOOS=$os GOARCH=$arch_ go build -o $bin -ldflags "-X main.version=$(git describe)" github.com/yuntan/gyazauto/cmd/gyazauto
}

if [ $# == 2 ]; then
  build_and_compress $1 $2
  exit 0
fi

for os in linux darwin windows; do
  build_and_compress $os 64
done

for os in linux windows; do
  build_and_compress $os 32
done
